package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines category model.
type Categories struct {
	ID             uuid.UUID `json:"id"       sql:"uuid"`
	Categories     string    `json:"categories" validate:"required" sql:"categories"`
	CreatedAt      time.Time `json:"createdat" sql:"createdat"`

}

// Query operations

// Gets category. Limit count and start position in db.
func GetCategories(db *sql.DB, start, count int) ([]Categories, error) {
	rows, err := db.Query(
		"SELECT id, categories, createdat FROM categories LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	category := []Categories{}

	// Store query results into category variable if no errors.
	for rows.Next() {
		var dt Categories
		if err := rows.Scan(&dt.ID, &dt.Categories, &dt.CreatedAt); err != nil {
			return nil, err
		}
		category = append(category, dt)
	}

	return category, nil
}


// Create new category and insert to database.
func (dt *Categories) CreateCategory(db *sql.DB) error {
	// Scan db after creation if user exists using new user id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO categories(categories, createdat) VALUES($1, $2) RETURNING id, categories, createdat", dt.Categories, timestamp).Scan(&dt.ID, &dt.Categories, &dt.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}


