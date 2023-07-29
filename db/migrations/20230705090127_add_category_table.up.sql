CREATE TABLE IF NOT EXISTS category(
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "subcategory" varchar,
    "brand" varchar NOT NULL,
    "rating" varchar NOT NULL,
    "store" INT NOT NULL,
    FOREIGN KEY ("store") REFERENCES "store" (id),
    "created_at" TIMESTAMP DEFAULT current_timestamp,
    "updated_at" TIMESTAMP DEFAULT current_timestamp
);
