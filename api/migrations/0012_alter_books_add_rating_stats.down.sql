DROP INDEX IF EXISTS idx_books_average_heart_score;
DROP INDEX IF EXISTS idx_books_average_poo_score;

ALTER TABLE books
DROP COLUMN average_heart_score,
DROP COLUMN average_poo_score,
DROP COLUMN total_ratings;