package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
)

// Defines book model.
type Books struct {
	ID               uuid.UUID `json:"id"       sql:"uuid"`
	BookID           uuid.UUID `json:"bookID" validate:"required" sql:"book_id"`
	NumberOfBooks    uint      `json:"numberOfBooks" validate:"required" sql:"number_of_books"`
	CreatedAt        time.Time `json:"createdAt" sql:"created_at"`
	DeletedAt        time.Time `json:"deletedAt" sql:"deleted_at"`
}

// Query operations

// Gets a specific book by id.
func (dt *Books) GetNumberBook(db *sql.DB) error {
	return db.QueryRow("SELECT book_id, number_of_book, created_at, deleted_at FROM books WHERE id=$1",
		dt.ID).Scan(&dt.BookID, &dt.NumberOfBooks,  &dt.CreatedAt, &dt.DeletedAt)
}

// Gets books. Limit count and start position in db.
func GetNumberBooks(db *sql.DB, field, sort string, limit, page int) ([]Books, error) {

	rows, err := db.Query(  "SELECT * FROM books ORDER BY $1 ,$2 LIMIT $3 OFFSET $4",
		field ,sort ,limit, limit*(page-1))

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	book := []Books{}

	// Store query results into book variable if no errors.
	for rows.Next() {
		var dt Books
		if err := rows.Scan(&dt.ID, &dt.BookID, &dt.NumberOfBooks, &dt.CreatedAt, &dt.DeletedAt); err != nil {
			return nil, err
		}
		book = append(book, dt)
	}

	return book, nil
}

// CRUD operations

// Create new  and insert to database.
func (dt *Books) CreateNumberBook(db *sql.DB) error {
	if dt.NumberOfBooks == 0 {
		return errors.New("books cannot be zero")
	}
	// Scan db after creation if book exists using new book id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO books(book_id, number_of_book, created_at, deleted_at) VALUES($1, $2, $3, $4) RETURNING id, book_id, number_of_book, created_at, deleted_at", dt.BookID, dt.NumberOfBooks, timestamp, timestamp).Scan(&dt.ID, &dt.BookID, &dt.NumberOfBooks, &dt.CreatedAt, &dt.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific book details by id.
func (dt *Books) UpdateNumberBook(db *sql.DB) error {
	if dt.NumberOfBooks == 0 {
		return errors.New("books cannot be zero")
	}
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE books SET book_id=$1, number_of_book=$2, deleted_at=$3 WHERE id=$4 RETURNING id,  book_id, number_of_book, created_at, deleted_at",dt.BookID, dt.NumberOfBooks, timestamp, dt.ID)

	return err
}

// Deletes a specific book by id.
func (dt *Books) DeleteAllBooks(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM books WHERE id=$1", dt.ID)

	return err
}
