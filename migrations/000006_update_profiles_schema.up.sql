-- migrations/006_update_profiles_schema.sql
-- Update user_profiles and artist_profiles tables

-- ============================================
-- UPDATE ARTIST_PROFILES TABLE
-- ============================================

-- Remove bio and monthly_listeners columns
ALTER TABLE artist_profiles DROP COLUMN IF EXISTS bio;
ALTER TABLE artist_profiles DROP COLUMN IF EXISTS monthly_listeners;

-- Add new columns to match user_profiles structure
ALTER TABLE artist_profiles ADD COLUMN IF NOT EXISTS full_name TEXT;
ALTER TABLE artist_profiles ADD COLUMN IF NOT EXISTS sex VARCHAR(10);
ALTER TABLE artist_profiles ADD COLUMN IF NOT EXISTS birthday DATE;

-- Rename image_url to avatar_url for consistency
ALTER TABLE artist_profiles RENAME COLUMN image_url TO avatar_url;
