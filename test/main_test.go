package test

import (
	"github.com/library/app"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/library/db"
	"github.com/spf13/viper"
)

// References App struct in app.go.
var a app.App

// References DB struct in app.go.
var d db.DB

//Generate new uuid for test
var testID = uuid.NewString()

// Executes before all other tests.
func TestMain(m *testing.M) {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	db_user := viper.GetString("TEST_DB_USERNAME")
	db_pass := viper.GetString("TEST_DB_PASSWORD")
	db_host := viper.GetString("TEST_DB_HOST")
	db_name := viper.GetString("TEST_DB_NAME")

	d.Initialize(db_user, db_pass, db_host, db_name)
	a.Initialize()
	ensureTableExists()
	// Executes tests.
	code := m.Run()
	// Cleans testing table is cleaned from database.
	clearTable()
	os.Exit(code)
}

// Helpers

// Executes http request using the router and returns response.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// Compares actual response to expected response of http request.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Ensures table needed for testing exists.
func ensureTableExists() {
	if _, err := d.Database.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// Clean test tables.
func clearTable() {
	d.Database.Exec("DELETE FROM admins")
	d.Database.Exec("DELETE FROM users")
	d.Database.Exec("DELETE FROM categories")
	d.Database.Exec("DELETE FROM authors")
	d.Database.Exec("DELETE FROM book")
	d.Database.Exec("DELETE FROM books")
	d.Database.Exec("DELETE FROM issue")
	d.Database.Exec("DELETE FROM acceptance")
}

// SQL query to create table.
const tableCreationQuery = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS admins (
		id uuid DEFAULT uuid_generate_v4 () unique,
		email varchar(225) NOT NULL UNIQUE,
		password varchar(225) NOT NULL,
		created_at timestamp NOT NULL,
		updated_at timestamp NOT NULL,
		primary key (id)
	);
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
    CREATE TABLE IF NOT EXISTS categories (
		id uuid DEFAULT uuid_generate_v4 () unique,
		name varchar(225) NOT NULL UNIQUE,
		created_at timestamp NOT NULL,
		primary key (id)
	);
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
	CREATE TABLE IF NOT EXISTS issue (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    user_id uuid,
	    books_id uuid,
	    return_date timestamp NOT NULL,
	    preliminary_cost float NOT NULL,
		created_at timestamp NOT NULL,
		primary key (id)
	);
	CREATE TABLE IF NOT EXISTS acceptance (
		id uuid DEFAULT uuid_generate_v4 () unique,
	    user_id uuid,
	    books_id uuid,
	    book_condition varchar(225) NOT NULL,
	    rating int NOT NULL,
	    discount float NOT NULL,
	    final_cost float NOT NULL,
	    photo varchar(225) NOT NULL,
	    created_at timestamp NOT NULL,
	    updated_at timestamp NOT NULL,
		primary key (id)
	);
`
