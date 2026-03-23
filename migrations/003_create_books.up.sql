CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    isbn TEXT UNIQUE,
    genre TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);