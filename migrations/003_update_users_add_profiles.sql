-- migrations/003_update_users_add_profiles.sql
-- Update users table and add user_profiles table

-- ============================================
-- UPDATE USERS TABLE
-- ============================================

-- Rename 'name' to 'username'
ALTER TABLE users RENAME COLUMN name TO username;

-- Rename 'password_hash' to 'password' (để khớp với struct)
ALTER TABLE users RENAME COLUMN password_hash TO password;

-- Drop columns không cần trong User struct
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS is_premium;

-- ============================================
-- CREATE USER_PROFILES TABLE
-- ============================================

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    show_profile BOOLEAN DEFAULT TRUE,
    full_name VARCHAR(100),
    avatar_url TEXT,
    sex VARCHAR(10),
    birthday DATE,
    country VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);
