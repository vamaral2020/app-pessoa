package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	//"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PessoasCreatedRequest struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

type PessoaResponse struct {
	ID       int64     `json:"id"`
	Nome     string    `json:"nome"`
	Email    string    `json:"email"`
	CriadoEm time.Time `json:"criadoEm"`
}

func main() {

	dsn := getenv("DB_DSN", "postgres://app:app@localhost:5432/appdb?sslmode=disable")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Erro abrindo conexao", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Erro ping no banco", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /pessoas", func(w http.ResponseWriter, r *http.Request) {
		handleCreatePessoa(w, r, db)
	})
	port := getenv("PORT", "8080")
	log.Println("API rodando em http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, logRequest(mux)); err != nil {
		log.Fatal(err)
	}

}
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
func handleCreatePessoa(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req PessoasCreatedRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "JSON invalido",
		})
		return
	}

	req.Nome = strings.TrimSpace(req.Nome)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Nome == "" || req.Email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "Nome e emails são obrigatorios",
		})
		return
	}

	if !strings.Contains(req.Email, "@") {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "email invalido",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //isso é um retorno de contexto com cancelamento
	defer cancel()

	var resp PessoaResponse
	err := db.QueryRowContext(ctx, `
		INSERT INTO pessoas (nome, email)
		VALUES ($1, $2)
		RETURNING id, nome, email, criado_em`, req.Nome, req.Email).Scan(&resp.ID, &resp.Nome, &resp.Email, &resp.CriadoEm)

	if err != nil {
		//tratamento do erro
		if isUniqueViolation(err) {
			writeJSON(w, http.StatusConflict, map[string]any{
				"error": "Email ja cadastrado",
			})
			return
		}

		log.Println("Erro ao inserir o pessoa", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Falha ao salvar pessoa",
		})
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique") ||
		errors.Is(err, sql.ErrNoRows) == false && strings.Contains(msg, "pessas_email_key")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func getenv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}
