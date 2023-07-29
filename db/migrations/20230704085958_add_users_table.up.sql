CREATE TABLE IF NOT EXISTS users (
  "id" bigserial PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "campus" varchar NOT NULL,
  "phone" varchar NOT NULL UNIQUE,
  "password" varchar NOT NULL,
  "usertype" varchar NOT NULL,
  "twofa" varchar NOT NULL,
  "active" varchar NOT NULL,
  "code" varchar NOT NULL,
  "wallet" varchar NOT NULL,
  "created_at" TIMESTAMP DEFAULT current_timestamp,
  "updated_at" TIMESTAMP DEFAULT current_timestamp
);
