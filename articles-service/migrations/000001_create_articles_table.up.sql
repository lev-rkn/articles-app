CREATE TABLE "articles" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "title" varchar(140) NOT NULL,
  "description" varchar(1000) NOT NULL,
  "photos" varchar[] NOT NULL,
  "timestamp" timestamptz NOT NULL DEFAULT (now())
);