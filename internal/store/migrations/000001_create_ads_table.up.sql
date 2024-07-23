CREATE TABLE "advertisements" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" text NOT NULL,
  "price" bigint NOT NULL,
  "photos" varchar[] NOT NULL
);
