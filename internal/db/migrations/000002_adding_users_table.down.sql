-- Drop indexes
DROP INDEX IF EXISTS idx_mailboxes_user_id;
DROP INDEX IF EXISTS idx_domains_user_id;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Drop foreign key columns
ALTER TABLE "domains" DROP COLUMN IF EXISTS "user_id";
ALTER TABLE "mailboxes" DROP COLUMN IF EXISTS "user_id";

-- Drop users table
DROP TABLE IF EXISTS "users";
