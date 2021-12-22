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
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);
`

// Schema for user table.
const USER_SCHEMA = `
	CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    firstname varchar(225) NOT NULL,
		surname varchar(225) NOT NULL,
	    secondname varchar(225) NOT NULL,
	    passport varchar(225) NOT NULL,
	    dateofbirth varchar(225) NOT NULL,
		email varchar(225) NOT NULL,
	    address varchar(225) NOT NULL,
	    indebtedness varchar(225) NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);
`
// Schema for category table.
const CATEGORY_SCHEMA = `
	CREATE TABLE IF NOT EXISTS categories (
		id uuid DEFAULT uuid_generate_v4 () unique,
		categories varchar(225) NOT NULL UNIQUE,
		createdat timestamp NOT NULL,
		primary key (id)
	);
`

// Schema for author table.
const AUTHOR_SCHEMA = `
	CREATE TABLE IF NOT EXISTS authors (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    firstname varchar(225) NOT NULL,
		surname varchar(225) NOT NULL,
	    dateofbirth varchar(225) NOT NULL,
		photo varchar(225) NOT NULL,
		createdat timestamp NOT NULL,
	    updatedat timestamp NOT NULL,
		primary key (id)
	);
`
// Schema for user table.
const BOOK_SCHEMA = `
	CREATE TABLE IF NOT EXISTS books (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    bookname varchar(225) NOT NULL,
		categoryID varchar(225) NOT NULL,
	    authorID varchar(225) NOT NULL,
	    cost float NOT NULL,
	    numberofbook int NOT NULL,
		photo varchar(225) NOT NULL,
	    priceperday float NOT NULL,
	    yearofpublishing int NOT NULL,
	    numberofpages int NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);
`
// Schema for lending table.
const LENDING_SCHEMA = `
	CREATE TABLE IF NOT EXISTS issue (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    userID varchar(225) NOT NULL,
	    bookID varchar(225) NOT NULL,
	    returndate varchar(225) NOT NULL,
	    preliminarycost float NOT NULL,
		createdat timestamp NOT NULL,
		primary key (id)
	);
`
// Schema for acceptance table.
const ACCEPTANCE_SCHEMA = `
	CREATE TABLE IF NOT EXISTS acceptance (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    userID varchar(225) NOT NULL,
	    bookID varchar(225) NOT NULL,
	    bookcondition varchar(225) NOT NULL,
	    rating int NOT NULL,
	    finalcost float NOT NULL,
	    photo varchar(225) NOT NULL,
	    createdat timestamp NOT NULL,
	    updatedat timestamp NOT NULL,
		primary key (id)
	);
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
	db.Database.Exec(LENDING_SCHEMA)
	db.Database.Exec(ACCEPTANCE_SCHEMA)
}
