CREATE TABLE "advertisements" (
  "id" bigserial PRIMARY KEY,
  "title" varchar(200) NOT NULL,
  "description" varchar(1000) NOT NULL,
  "price" decimal NOT NULL,
  "photos" varchar[] NOT NULL,
  "timestamp" timestamptz NOT NULL DEFAULT (now())
);
