CREATE TABLE "refresh_sessions" (
    "id" bigserial PRIMARY KEY,
    "refresh_token" varchar(255) NOT NULL,
    "fingerprint" varchar NOT NULL,
    "user_email" varchar(255) NOT NULL,
    "app_id" integer NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "refresh_sessions" ADD FOREIGN KEY ("user_email") REFERENCES "users" ("email");
ALTER TABLE "refresh_sessions" ADD FOREIGN KEY ("app_id") REFERENCES "apps" ("id");

CREATE INDEX idx_refresh_token ON refresh_sessions (refresh_token);
