-- migrations/004_remove_role.sql
-- Remove role column from users table

-- Drop index for role queries
DROP INDEX IF EXISTS idx_users_role;

-- Remove role column
ALTER TABLE users DROP COLUMN IF EXISTS role;
