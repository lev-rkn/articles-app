CREATE TABLE "apps" (
    "id" serial PRIMARY KEY,
    "name" varchar NOT NULL UNIQUE,
    "secret" text NOT NULL UNIQUE
);
