-- Rollback 005_add_artist_profiles
DROP INDEX IF EXISTS idx_artist_profiles_artist_id;
DROP TABLE IF EXISTS artist_profiles CASCADE;
ALTER TABLE artists ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE artists ADD COLUMN IF NOT EXISTS image_url TEXT;
ALTER TABLE artists ADD COLUMN IF NOT EXISTS monthly_listeners INT DEFAULT 0;
ALTER TABLE artists DROP COLUMN IF EXISTS updated_at;
