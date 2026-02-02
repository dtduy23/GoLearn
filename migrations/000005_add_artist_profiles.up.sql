-- migrations/005_add_artist_profiles.sql
-- Update artists table and add artist_profiles table

-- ============================================
-- UPDATE ARTISTS TABLE
-- ============================================

-- Add updated_at column if not exists
ALTER TABLE artists ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Move profile fields to artist_profiles
-- First, drop the columns that will be moved to profile
ALTER TABLE artists DROP COLUMN IF EXISTS bio;
ALTER TABLE artists DROP COLUMN IF EXISTS image_url;  
ALTER TABLE artists DROP COLUMN IF EXISTS monthly_listeners;

-- ============================================
-- CREATE ARTIST_PROFILES TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS artist_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    artist_id UUID UNIQUE REFERENCES artists(id) ON DELETE CASCADE,
    bio TEXT,
    image_url TEXT,
    monthly_listeners INT DEFAULT 0,
    country VARCHAR(100),
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_artist_profiles_artist_id ON artist_profiles(artist_id);
