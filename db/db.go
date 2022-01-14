package db

import (
	"database/sql"
	"fmt"
	"log"
)

type DB struct {
	Database *sql.DB
}

const DB_SETUP = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

// Schema for admin table.
const ADMIN_SCHEMA = `
	CREATE TABLE IF NOT EXISTS admins (
		id uuid DEFAULT uuid_generate_v4 () unique,
		email varchar(225) NOT NULL UNIQUE,
		password varchar(225) NOT NULL,
		created_at timestamp NOT NULL,
		updated_at timestamp NOT NULL,
		primary key (id)
	);
`

// Schema for user table.
const USER_SCHEMA = `
	CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    firstname varchar(225) NOT NULL,
		surname varchar(225) NOT NULL,
	    second_name varchar(225) NOT NULL,
	    passport varchar(225) NOT NULL,
	    date_of_birth varchar(225) NOT NULL,
		email varchar(225) NOT NULL,
	    address varchar(225) NOT NULL,
	    indebtedness varchar(225) NOT NULL,
		created_at timestamp NOT NULL,
		updated_at timestamp NOT NULL,
		primary key (id)
	);
`
// Schema for category table.
const CATEGORY_SCHEMA = `
	CREATE TABLE IF NOT EXISTS categories (
		id uuid DEFAULT uuid_generate_v4 () unique,
		name varchar(225) NOT NULL UNIQUE,
		created_at timestamp NOT NULL,
		primary key (id)
	);
`

// Schema for author table.
const AUTHOR_SCHEMA = `
	CREATE TABLE IF NOT EXISTS authors (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    firstname varchar(225) NOT NULL,
		surname varchar(225) NOT NULL,
	    date_of_birth varchar(225) NOT NULL,
		photo varchar(225) NOT NULL,
		created_at timestamp NOT NULL,
	    updated_at timestamp NOT NULL,
		primary key (id)
	);
`
// Schema for book table.
const BOOK_SCHEMA = `
	CREATE TABLE IF NOT EXISTS book (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    name varchar(225) NOT NULL,
		category_id uuid,
	    author_id uuid,
	    cost float NOT NULL,
	    price_per_day float NOT NULL,
		photo varchar(225) NOT NULL,
	    year_of_publishing int NOT NULL,
	    number_of_pages int NOT NULL,
		created_at timestamp NOT NULL,
		updated_at timestamp NOT NULL,
		primary key (id)
	);

ALTER TABLE book
    ADD CONSTRAINT fk_categories_book
        FOREIGN KEY (category_id)
            REFERENCES categories(id)
            ON DELETE CASCADE;

ALTER TABLE book
    ADD CONSTRAINT fk_authors_book
        FOREIGN KEY (author_id)
            REFERENCES authors(id)
            ON DELETE CASCADE;
`
// Schema for books table.
const BOOKS_SCHEMA = `
	CREATE TABLE IF NOT EXISTS books (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    book_id uuid,
		category_id uuid,
	    author_id uuid,
	    number_of_book int NOT NULL,
		created_at timestamp NOT NULL,
		deleted_at timestamp NOT NULL,
		primary key (id)
	);

ALTER TABLE books
    ADD CONSTRAINT fk_book_books
        FOREIGN KEY (book_id)
            REFERENCES book(id)
            ON DELETE CASCADE;

ALTER TABLE books
    ADD CONSTRAINT fk_categories_books
        FOREIGN KEY (category_id)
            REFERENCES categories(id)
            ON DELETE CASCADE;

ALTER TABLE books
    ADD CONSTRAINT fk_authors_books
        FOREIGN KEY (author_id)
            REFERENCES authors(id)
            ON DELETE CASCADE;
`
// Schema for lending table.
const ISSUE_SCHEMA = `
	CREATE TABLE IF NOT EXISTS issue (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    user_id uuid,
	    book_id uuid,
	    return_date timestamp NOT NULL,
	    preliminary_cost float NOT NULL,
		created_at timestamp NOT NULL,
		primary key (id)
	);

ALTER TABLE issue
    ADD CONSTRAINT fk_users_issue
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE;

ALTER TABLE issue
    ADD CONSTRAINT fk_book_issue
        FOREIGN KEY (book_id)
            REFERENCES book(id)
            ON DELETE CASCADE;
`
// Schema for acceptance table.
const ACCEPTANCE_SCHEMA = `
	CREATE TABLE IF NOT EXISTS acceptance (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    user_id uuid,
	    book_id uuid,
	    book_condition varchar(225) NOT NULL,
	    discount float NOT NULL,
	    final_cost float NOT NULL,
	    photo varchar(225) NOT NULL,
	    created_at timestamp NOT NULL,
	    updated_at timestamp NOT NULL,
		primary key (id)
	);

ALTER TABLE acceptance
    ADD CONSTRAINT fk_users_acceptance
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE;

ALTER TABLE acceptance
    ADD CONSTRAINT fk_book_acceptance
        FOREIGN KEY (book_id)
            REFERENCES book(id)
            ON DELETE CASCADE;
`
const POPULAR_BOOKS_SCHEMA = `
	CREATE TABLE IF NOT EXISTS popular (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    book_id uuid,
	    rating int NOT NULL,
	    views int NOT NULL,
	    created_at timestamp NOT NULL,
	    updated_at timestamp NOT NULL,
		primary key (id)
	);

ALTER TABLE popular
    ADD CONSTRAINT fk_book_popular
        FOREIGN KEY (book_id)
            REFERENCES book(id)
            ON DELETE CASCADE;
`
// Receives database credentials and connects to database.
func (db *DB) Initialize(user, password, dbhost, dbname string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, dbhost, dbname)

	var err error
	db.Database, err = sql.Open("postgres", connectionString)
	// Log errors.
	if err != nil {
		log.Fatal(err)
	}
	db.Database.Exec(DB_SETUP)
	db.Database.Exec(ADMIN_SCHEMA)
	db.Database.Exec(USER_SCHEMA)
    db.Database.Exec(CATEGORY_SCHEMA)
	db.Database.Exec(AUTHOR_SCHEMA)
	db.Database.Exec(BOOK_SCHEMA)
	db.Database.Exec(BOOKS_SCHEMA)
	db.Database.Exec(ISSUE_SCHEMA)
	db.Database.Exec(ACCEPTANCE_SCHEMA)
	db.Database.Exec(POPULAR_BOOKS_SCHEMA)
}
