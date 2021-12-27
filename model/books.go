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
	CategoryID       uuid.UUID `json:"categoryID" validate:"required" sql:"category_id"`
	AuthorID         uuid.UUID `json:"authorID" validate:"required" sql:"author_id"`
	NumberOfBooks    uint      `json:"numberOfBooks" validate:"required" sql:"number_of_books"`
	CreatedAt        time.Time `json:"createdAt" sql:"created_at"`
	DeletedAt        time.Time `json:"deletedAt" sql:"deleted_at"`
}

// Query operations

// Gets a specific book by id.
func (dt *Books) GetNumberBook(db *sql.DB) error {
	return db.QueryRow("SELECT book_id, category_id, author_id, number_of_book, created_at, deleted_at FROM books WHERE id=$1",
		dt.ID).Scan(&dt.BookID,  &dt.CategoryID, &dt.AuthorID, &dt.NumberOfBooks,  &dt.CreatedAt, &dt.DeletedAt)
}

// Gets books. Limit count and start position in db.
func GetNumberBooks(db *sql.DB, start, count int) ([]Books, error) {
	rows, err := db.Query(
		"SELECT id, book_id, category_id, author_id, number_of_book, created_at, deleted_at FROM books LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	book := []Books{}

	// Store query results into book variable if no errors.
	for rows.Next() {
		var dt Books
		if err := rows.Scan(&dt.ID, &dt.BookID, &dt.CategoryID, &dt.AuthorID, &dt.NumberOfBooks, &dt.CreatedAt, &dt.DeletedAt); err != nil {
			return nil, err
		}
		book = append(book, dt)
	}

	return book, nil
}

// CRUD operations

// Create new  and insert to database.
func (dt *Books) CreateNumberBook(db *sql.DB) error {
	// Scan db after creation if book exists using new book id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO books(book_id, category_id, author_id, number_of_book, created_at, deleted_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, book_id, category_id, author_id, number_of_book, created_at, deleted_at", dt.BookID, dt.CategoryID, dt.AuthorID, dt.NumberOfBooks, timestamp, timestamp).Scan(&dt.ID, &dt.BookID, &dt.CategoryID, &dt.AuthorID, &dt.NumberOfBooks, &dt.CreatedAt, &dt.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific book details by id.
func (dt *Books) UpdateNumberBook(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE books SET book_id=$1, category_id=$2, author_id=$3, number_of_book=$4, deleted_at=$5 WHERE id=$6 RETURNING id,  book_id, category_id, author_id, number_of_book, created_at, deleted_at",dt.BookID, dt.CategoryID, dt.AuthorID, dt.NumberOfBooks, timestamp, dt.ID)

	return err
}

// Deletes a specific book by id.
func (dt *Books) DeleteAllBooks(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM books WHERE id=$1", dt.ID)

	return err
}
func (r *Books) Restrictions() error {
	if r.NumberOfBooks > 5 {
		return errors.New("issuing more than five books is prohibited")
	}
	return nil
}

func (books *Books) Validate() error {
	if books.NumberOfBooks == 0 {
		return errors.New("books cannot be zero")
	}
	return nil
}