CREATE TABLE "users" (
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar,
  "provider" varchar,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "documents" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "filename" varchar,
  "file_type" varchar,
  "uploaded_at" timestamp DEFAULT (now())
);

CREATE TABLE "extracted_texts" (
  "id" varchar PRIMARY KEY,
  "document_id" varchar NOT NULL,
  "content" text,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "documents" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("username");

ALTER TABLE "extracted_texts" ADD FOREIGN KEY ("document_id") REFERENCES "documents" ("id");
