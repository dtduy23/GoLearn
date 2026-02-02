-- Rollback 003_update_users_add_profiles
DROP INDEX IF EXISTS idx_users_role;
ALTER TABLE users DROP COLUMN IF EXISTS role;
DROP INDEX IF EXISTS idx_user_profiles_user_id;
DROP TABLE IF EXISTS user_profiles CASCADE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_premium BOOLEAN DEFAULT FALSE;
ALTER TABLE users RENAME COLUMN password TO password_hash;
ALTER TABLE users RENAME COLUMN username TO name;
