CREATE TABLE  IF NOT EXISTS store (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "description" varchar,
    "address" varchar,
    "review" varchar,
    "follower" varchar,
    "created_at" TIMESTAMP DEFAULT current_timestamp,
    "updated_at" TIMESTAMP DEFAULT current_timestamp,
    "owner" INT NOT NULL,
    FOREIGN KEY ("owner") REFERENCES "users" (id)
);
