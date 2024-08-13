CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "email" varchar(255) NOT NULL UNIQUE,
    "pass_hash" BYTEA NOT NULL
);

CREATE INDEX idx_email ON users (email);
