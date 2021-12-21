package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines book model.
type Book struct {
	ID               uuid.UUID `json:"id"       sql:"uuid"`
	BookName         string    `json:"bookname" validate:"required" sql:"bookname"`
	CategoryID       string    `json:"categoryID" validate:"required" sql:"categoryID"`
	AuthorID         string    `json:"authorID" validate:"required" sql:"authorID"`
	Cost             float32   `json:"cost" validate:"required" sql:"cost"`
	NumberOfBook     uint      `json:"numberofbook" validate:"required" sql:"numberofbook"`
	Photo            string    `json:"photo" validate:"required" sql:"photo"`
	PricePerDay      float32   `json:"priceperday" validate:"required" sql:"priceperday"`
	YearOfPublishing uint      `json:"yearofpublishing" validate:"required" sql:"yearofpublishing"`
	NumberOfPages    uint      `json:"numberofpages" validate:"required" sql:"numberofpages"`
	CreatedAt        time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt        time.Time `json:"updatedat" sql:"updatedat"`
}

// Query operations

// Gets a specific book by id.
func (dt *Book) GetBook(db *sql.DB) error {
	return db.QueryRow("SELECT bookname, categoryID, authorID, cost, numberofbook, photo, priceperday, yearofpublishing, numberofpages, createdat, updatedat FROM books WHERE id=$1",
		dt.ID).Scan(&dt.BookName, &dt.CategoryID, &dt.AuthorID, &dt.Cost, &dt.NumberOfBook, &dt.Photo, &dt.PricePerDay, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets books. Limit count and start position in db.
func GetBooks(db *sql.DB, start, count int) ([]Book, error) {
	rows, err := db.Query(
		"SELECT id, bookname, categoryID, authorID, cost, numberofbook, photo, priceperday, yearofpublishing, numberofpages, createdat, updatedat FROM books LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	book := []Book{}

	// Store query results into book variable if no errors.
	for rows.Next() {
		var dt Book
		if err := rows.Scan(&dt.ID, &dt.BookName, &dt.CategoryID, &dt.AuthorID, &dt.Cost, &dt.NumberOfBook, &dt.Photo, &dt.PricePerDay, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		book = append(book, dt)
	}

	return book, nil
}

// CRUD operations

// Create new user and insert to database.
func (dt *Book) CreateBook(db *sql.DB) error {
	// Scan db after creation if book exists using new book id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO books(bookname, categoryID, authorID, cost, numberofbook, photo, priceperday, yearofpublishing, numberofpages, createdat, updatedat) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, bookname, categoryID, authorID, cost, numberofbook, photo, priceperday, yearofpublishing, numberofpages, createdat, updatedat", dt.BookName, dt.CategoryID, dt.AuthorID, dt.Cost, dt.NumberOfBook, dt.Photo, dt.PricePerDay, dt.YearOfPublishing, dt.NumberOfPages, timestamp, timestamp).Scan(&dt.ID, &dt.BookName, &dt.CategoryID, &dt.AuthorID, &dt.Cost, &dt.NumberOfBook, &dt.Photo, &dt.PricePerDay, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific book details by id.
func (dt *Book) UpdateBook(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE books SET bookname=$1, categoryID=$2, authorID=$3, cost=$4, numberofbook=$5, photo=$6, priceperday=$7, yearofpublishing=$8, numberofpages=$9, updatedat=$10 WHERE id=$11 RETURNING id, bookname, categoryID, authorID, cost, numberofbook, photo, priceperday, yearofpublishing, numberofpages, createdat, updatedat", dt.BookName, dt.CategoryID, dt.AuthorID, dt.Cost, dt.NumberOfBook, dt.Photo, dt.PricePerDay, dt.YearOfPublishing, dt.NumberOfPages, timestamp, dt.ID)

	return err
}

// Deletes a specific book by id.
func (dt *Book) DeleteBook(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM books WHERE id=$1", dt.ID)

	return err
}