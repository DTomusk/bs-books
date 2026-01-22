DROP TABLE book_author;

ALTER TABLE books ADD COLUMN author_id UUID REFERENCES authors(id);