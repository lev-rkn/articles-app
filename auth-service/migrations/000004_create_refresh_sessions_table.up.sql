CREATE TABLE "refreshSessions" (
    "id" bigserial PRIMARY KEY,
    "refresh_token" varchar(255) NOT NULL,
    "fingerprint" varchar NOT NULL,
    "user_email" bigint NOT NULL,
    "app_id" int NOT NULL
    "createdAt" timestamptz NOT NULL DEFAULT now()
);
ALTER TABLE "refreshSessions" ADD FOREIGN KEY ("user_email") REFERENCES "users" ("email");
ALTER TABLE "refreshSessions" ADD FOREIGN KEY ("app_id") REFERENCES "apps" ("id");

CREATE INDEX idx_refresh_token ON refreshSessions (refresh_token);
