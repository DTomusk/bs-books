ALTER TABLE books DROP COLUMN IF EXISTS cover_img_url;
ALTER TABLE books DROP COLUMN IF EXISTS synopsis;

DROP INDEX IF EXISTS idx_books_title_trgm;