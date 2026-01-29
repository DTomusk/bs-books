CREATE TABLE author_alias (
    author_id UUID REFERENCES authors(id),
    alias TEXT NOT NULL,
    PRIMARY KEY (author_id, alias)
);