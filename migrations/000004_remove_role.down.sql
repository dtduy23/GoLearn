-- Rollback 004_remove_role
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user' NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
