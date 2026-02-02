ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);
-- Tạo function
CREATE OR REPLACE FUNCTION prevent_username_update()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.username IS DISTINCT FROM NEW.username THEN
        RAISE EXCEPTION 'username cannot be changed';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Tạo trigger
CREATE TRIGGER username_immutable
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION prevent_username_update();