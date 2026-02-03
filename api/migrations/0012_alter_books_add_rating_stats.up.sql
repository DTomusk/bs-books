ALTER TABLE books
ADD COLUMN average_heart_score FLOAT8 NOT NULL DEFAULT 0,
ADD COLUMN average_poo_score FLOAT8 NOT NULL DEFAULT 0,
ADD COLUMN total_ratings INT NOT NULL DEFAULT 0;

CREATE INDEX idx_books_average_heart_score
ON books (average_heart_score);

CREATE INDEX idx_books_average_poo_score
ON books (average_poo_score);