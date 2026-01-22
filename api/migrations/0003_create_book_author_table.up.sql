ALTER TABLE books DROP COLUMN author_id;

CREATE TABLE book_author (
    book_id UUID REFERENCES books(id),
    author_id UUID REFERENCES authors(id),
    PRIMARY KEY (book_id, author_id)
);