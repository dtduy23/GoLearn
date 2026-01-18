-- migrations/002_seed_data.sql
-- Sample data for testing

-- Insert sample genres
INSERT INTO genres (name, image_url) VALUES
    ('Pop', '/static/genres/pop.jpg'),
    ('Rock', '/static/genres/rock.jpg'),
    ('Hip Hop', '/static/genres/hiphop.jpg'),
    ('R&B', '/static/genres/rnb.jpg'),
    ('Electronic', '/static/genres/electronic.jpg'),
    ('Jazz', '/static/genres/jazz.jpg'),
    ('Classical', '/static/genres/classical.jpg'),
    ('Country', '/static/genres/country.jpg');

-- Insert sample artists
INSERT INTO artists (id, name, bio, image_url) VALUES
    ('a1111111-1111-1111-1111-111111111111', 'The Weeknd', 'Canadian singer and songwriter', '/static/artists/weeknd.jpg'),
    ('a2222222-2222-2222-2222-222222222222', 'Taylor Swift', 'American singer-songwriter', '/static/artists/taylor.jpg'),
    ('a3333333-3333-3333-3333-333333333333', 'Ed Sheeran', 'English singer-songwriter', '/static/artists/ed.jpg'),
    ('a4444444-4444-4444-4444-444444444444', 'Dua Lipa', 'English singer', '/static/artists/dua.jpg');

-- Insert sample albums
INSERT INTO albums (id, artist_id, title, cover_url, release_date, album_type) VALUES
    ('b1111111-1111-1111-1111-111111111111', 'a1111111-1111-1111-1111-111111111111', 'After Hours', '/static/albums/after_hours.jpg', '2020-03-20', 'album'),
    ('b2222222-2222-2222-2222-222222222222', 'a2222222-2222-2222-2222-222222222222', '1989', '/static/albums/1989.jpg', '2014-10-27', 'album'),
    ('b3333333-3333-3333-3333-333333333333', 'a3333333-3333-3333-3333-333333333333', 'Divide', '/static/albums/divide.jpg', '2017-03-03', 'album'),
    ('b4444444-4444-4444-4444-444444444444', 'a4444444-4444-4444-4444-444444444444', 'Future Nostalgia', '/static/albums/future_nostalgia.jpg', '2020-03-27', 'album');

-- Insert sample songs
INSERT INTO songs (id, album_id, title, duration, file_url, track_number) VALUES
    ('c1111111-1111-1111-1111-111111111111', 'b1111111-1111-1111-1111-111111111111', 'Blinding Lights', 200, '/static/music/blinding_lights.mp3', 1),
    ('c1111111-1111-1111-1111-111111111112', 'b1111111-1111-1111-1111-111111111111', 'Save Your Tears', 215, '/static/music/save_your_tears.mp3', 2),
    ('c2222222-2222-2222-2222-222222222221', 'b2222222-2222-2222-2222-222222222222', 'Shake It Off', 219, '/static/music/shake_it_off.mp3', 1),
    ('c2222222-2222-2222-2222-222222222222', 'b2222222-2222-2222-2222-222222222222', 'Blank Space', 231, '/static/music/blank_space.mp3', 2),
    ('c3333333-3333-3333-3333-333333333331', 'b3333333-3333-3333-3333-333333333333', 'Shape of You', 234, '/static/music/shape_of_you.mp3', 1),
    ('c3333333-3333-3333-3333-333333333332', 'b3333333-3333-3333-3333-333333333333', 'Perfect', 263, '/static/music/perfect.mp3', 2),
    ('c4444444-4444-4444-4444-444444444441', 'b4444444-4444-4444-4444-444444444444', 'Dont Start Now', 183, '/static/music/dont_start_now.mp3', 1),
    ('c4444444-4444-4444-4444-444444444442', 'b4444444-4444-4444-4444-444444444444', 'Levitating', 203, '/static/music/levitating.mp3', 2);

-- Link songs to artists
INSERT INTO song_artists (song_id, artist_id, is_primary) VALUES
    ('c1111111-1111-1111-1111-111111111111', 'a1111111-1111-1111-1111-111111111111', true),
    ('c1111111-1111-1111-1111-111111111112', 'a1111111-1111-1111-1111-111111111111', true),
    ('c2222222-2222-2222-2222-222222222221', 'a2222222-2222-2222-2222-222222222222', true),
    ('c2222222-2222-2222-2222-222222222222', 'a2222222-2222-2222-2222-222222222222', true),
    ('c3333333-3333-3333-3333-333333333331', 'a3333333-3333-3333-3333-333333333333', true),
    ('c3333333-3333-3333-3333-333333333332', 'a3333333-3333-3333-3333-333333333333', true),
    ('c4444444-4444-4444-4444-444444444441', 'a4444444-4444-4444-4444-444444444444', true),
    ('c4444444-4444-4444-4444-444444444442', 'a4444444-4444-4444-4444-444444444444', true);

-- Insert a test user (password: "password123")
INSERT INTO users (id, email, password_hash, name) VALUES
    ('u1111111-1111-1111-1111-111111111111', 'test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjqSw5Y.O8TFUBxMTpJyTGV0lVpRQa', 'Test User');
