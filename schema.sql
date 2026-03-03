CREATE TABLE IF NOT EXISTS pessoas (
                                       id BIGSERIAL PRIMARY KEY,
                                       nome TEXT NOT NULL,
                                       email TEXT NOT NULL UNIQUE,
                                       criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );