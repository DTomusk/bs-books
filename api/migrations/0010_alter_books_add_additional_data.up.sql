ALTER TABLE books ADD COLUMN cover_img_url TEXT NULL;
ALTER TABLE books ADD COLUMN synopsis TEXT NULL;

CREATE INDEX idx_books_title_trgm
ON books
USING gin (title gin_trgm_ops);