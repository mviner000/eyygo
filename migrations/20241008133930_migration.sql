-- +migrate Up
CREATE TABLE IF NOT EXISTS borrowers (
 id integer PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL,
 email text NOT NULL,
 password text NOT NULL,
 created_at datetime,
 updated_at datetime
);

CREATE TABLE IF NOT EXISTS books (
 id integer PRIMARY KEY AUTOINCREMENT,
 title text NOT NULL,
 author_id integer NOT NULL,
 published_at datetime,
 category_id integer NOT NULL,
 stock integer NOT NULL
);

CREATE TABLE IF NOT EXISTS authors (
 id integer PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL
);

CREATE TABLE IF NOT EXISTS rentals (
 id integer PRIMARY KEY AUTOINCREMENT,
 borrower_id integer NOT NULL,
 book_id integer NOT NULL,
 rented_at datetime NOT NULL,
 due_at datetime NOT NULL,
 returned_at datetime
);

CREATE TABLE IF NOT EXISTS categories (
 id integer PRIMARY KEY AUTOINCREMENT,
 name text NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS borrowers;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS rentals;