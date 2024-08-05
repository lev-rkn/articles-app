CREATE TABLE "comments" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "article_id" bigint NOT NULL,
  "text" varchar(500) NOT NULL,
  "timestamp" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "comments" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id");