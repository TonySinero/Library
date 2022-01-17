package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
)

// Defines book model.
type Book struct {
	ID               uuid.UUID `json:"id"       sql:"uuid"`
	Name             string    `json:"name" validate:"required" sql:"name"`
	CategoryID       uuid.UUID `json:"categoryID" validate:"required" sql:"category_id"`
	AuthorID         uuid.UUID `json:"authorID" validate:"required" sql:"author_id"`
	Cost             float32   `json:"cost" validate:"required" sql:"cost"`
	PricePerDay      float32   `json:"pricePerDay" validate:"required" sql:"price_per_day"`
	Photo            string    `json:"photo" validate:"required" sql:"photo"`
	YearOfPublishing uint      `json:"yearOfPublishing" validate:"required" sql:"year_of_publishing"`
	NumberOfPages    uint      `json:"numberOfPages" validate:"required" sql:"number_of_pages"`
	Views            uint      `json:"views" validate:"required" sql:"views"`
	CreatedAt        time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations
func GetPopularBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query(
		"SELECT id, name, category_id, author_id, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at FROM book ORDER BY views DESC LIMIT 3",
	)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	book := []Book{}

	// Store query results into user variable if no errors.
	for rows.Next() {
		var dt Book
		if err := rows.Scan(&dt.ID, &dt.Name, &dt.CategoryID, &dt.AuthorID, &dt.Cost, &dt.PricePerDay, &dt.Photo, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		book = append(book, dt)
	}
	return book, nil
}
// Gets a specific book by name.
func (dt *Book) GetBook(db *sql.DB) error {
	return db.QueryRow("SELECT id, name, category_id, author_id, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at FROM book WHERE name=$1",
		dt.Name).Scan(&dt.ID, &dt.Name, &dt.CategoryID, &dt.AuthorID, &dt.Cost, &dt.PricePerDay, &dt.Photo, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets books. Limit count and start position in db.
func GetBooks(db *sql.DB, start, count int) ([]Book, error) {
	rows, err := db.Query(
		"SELECT id, name, category_id, author_id, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at FROM book ORDER BY name LIMIT $1 OFFSET $2",
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
		if err := rows.Scan(&dt.ID, &dt.Name, &dt.CategoryID, &dt.AuthorID, &dt.Cost, &dt.PricePerDay, &dt.Photo, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		book = append(book, dt)
	}

	return book, nil
}

// CRUD operations

// Create new book and insert to database.
func (dt *Book) CreateBook(db *sql.DB) error {
	// Scan db after creation if book exists using new book id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO book(name, category_id, author_id, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, name, category_id, author_id, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at", dt.Name, dt.CategoryID, dt.AuthorID, dt.Cost, dt.PricePerDay, dt.Photo, dt.YearOfPublishing, dt.NumberOfPages, dt.Views, timestamp, timestamp).Scan(&dt.ID, &dt.Name, &dt.CategoryID, &dt.AuthorID, &dt.Cost, &dt.PricePerDay, &dt.Photo, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific book details by id.
func (dt *Book) UpdateBook(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE book SET name=$1, category_id=$2, author_id=$3, cost=$4, price_per_day=$5, photo=$6, year_of_publishing=$7, number_of_pages=$8, views=$9, updated_at=$10 WHERE id=$11 RETURNING id, name, category_id, author_id, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at", dt.Name, dt.CategoryID, dt.AuthorID, dt.Cost, dt.PricePerDay, dt.Photo, dt.YearOfPublishing, dt.NumberOfPages, dt.Views, timestamp, dt.ID)

	return err
}

// Deletes a specific book by id.
func (dt *Book) DeleteBook(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM book WHERE id=$1", dt.ID)

	return err
}

func (dt *Book) Validate() error {
	if dt.Name == "" {
		return errors.New("name is required")
	}
	if dt.Cost == 0 {
		return errors.New("cost cannot be zero")
	}
	if dt.PricePerDay == 0 {
		return errors.New("pricePerDay is required")
	}
	if dt.Photo == "" {
		return errors.New("photo is required")
	}
	if dt.YearOfPublishing == 0 {
		return errors.New("yearOfPublishing cannot be zero")
	}
	if dt.NumberOfPages == 0 {
		return errors.New("numberOfPages cannot be zero")
	}
	return nil
}