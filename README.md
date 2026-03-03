# 🚀 app-pessoa

API simples para cadastro de pessoas, desenvolvida em Go, com foco em aprendizado, simplicidade e clareza estrutural.

---

## 📌 Sobre o Projeto

O **app-pessoa** é uma aplicação REST que permite:

- ✅ Cadastrar pessoas
- ✅ Listar pessoas
- ✅ Buscar pessoa por ID
- ✅ Persistência em banco de dados
- ✅ Estrutura enxuta e sem frameworks pesados

Projeto criado com objetivo de estudo da linguagem Go e seus conceitos fundamentais como:

- Organização por pacotes
- Manipulação de JSON
- Uso de banco de dados
- Tratamento de erros
- Boas práticas de estrutura modular

---

## 🛠️ Tecnologias Utilizadas

- **Go**
- **net/http**
- **encoding/json**
- **database/sql**
- Banco de dados (ex: PostgreSQL ou SQLite)
- Migrations (se aplicável)

---

## 📂 Estrutura do Projeto
```app-pessoa/
│
├── main.go
├── go.mod
├── handler/
│ └── pessoa_handler.go
├── service/
│ └── pessoa_service.go
├── repository/
│ └── pessoa_repository.go
├── model/
│ └── pessoa.go
└── database/
└── connection.go```


Arquitetura simples seguindo separação de responsabilidades:

- **handler** → camada HTTP
- **service** → regras de negócio
- **repository** → acesso a dados
- **model** → entidades
- **database** → conexão

---

## ⚙️ Como Executar

### 1️⃣ Clonar o repositório

```bash
git clone https://github.com/seu-usuario/app-pessoa.git
cd app-pessoa

go mod tidy

export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=usuario
export DB_PASSWORD=senha
export DB_NAME=app_pessoa

go run .

http://localhost:8080

POST /pessoas
Content-Type: application/json

{
  "nome": "Vagner",
  "email": "vagner@email.com"
}

{
  "id": 1,
  "nome": "Vagner",
  "email": "vagner@email.com",
  "criadoEm": "2026-03-03T00:32:22Z"
}```
