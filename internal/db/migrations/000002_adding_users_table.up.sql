CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "hashed_password" VARCHAR(255) NOT NULL,
    "full_name" VARCHAR(255),
    "role" VARCHAR(50) NOT NULL DEFAULT 'user',
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for search optimization
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);

-- Add relationships with existing tables
ALTER TABLE "domains" ADD COLUMN "user_id" INTEGER NOT NULL REFERENCES "users"(id) ON DELETE CASCADE;
ALTER TABLE "mailboxes" ADD COLUMN "user_id" INTEGER NOT NULL REFERENCES "users"(id) ON DELETE CASCADE;

-- Indexes for foreign keys
CREATE INDEX idx_domains_user_id ON domains(user_id);
CREATE INDEX idx_mailboxes_user_id ON mailboxes(user_id);
