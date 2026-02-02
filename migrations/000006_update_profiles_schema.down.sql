-- Rollback 006_update_profiles_schema
ALTER TABLE artist_profiles RENAME COLUMN avatar_url TO image_url;
ALTER TABLE artist_profiles DROP COLUMN IF EXISTS full_name;
ALTER TABLE artist_profiles DROP COLUMN IF EXISTS sex;
ALTER TABLE artist_profiles DROP COLUMN IF EXISTS birthday;
ALTER TABLE artist_profiles ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE artist_profiles ADD COLUMN IF NOT EXISTS monthly_listeners INT DEFAULT 0;
