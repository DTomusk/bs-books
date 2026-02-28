ALTER TABLE users
ADD COLUMN username VARCHAR(50);

UPDATE users 
SET username = 'user_' || LEFT(MD5(CAST(id AS TEXT)), 8)
WHERE username IS NULL;

ALTER TABLE users
ALTER COLUMN username SET NOT NULL;