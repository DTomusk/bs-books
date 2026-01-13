CREATE TABLE authors (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    author_id UUID REFERENCES authors(id)
);

CREATE TABLE users (
    id UUID PRIMARY KEY
);

CREATE TABLE ratings (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    book_id UUID REFERENCES books(id),
    heart_score NUMERIC(2,1) NOT NULL,
    poo_score NUMERIC(2,1) NOT NULL
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY,
    rating_id UUID REFERENCES ratings(id),
    review TEXT NOT NULL
);