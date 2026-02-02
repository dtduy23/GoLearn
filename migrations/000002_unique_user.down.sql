-- Rollback 002_unique_user
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_unique;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_unique;
