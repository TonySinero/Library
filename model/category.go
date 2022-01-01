package model

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// Defines category model.
type Categories struct {
	ID             uuid.UUID `json:"id"       sql:"uuid"`
	Name           string    `json:"name" validate:"required" sql:"name"`
	CreatedAt      time.Time `json:"createdAt" sql:"created_at"`

}

// Query operations

// Gets category. Limit count and start position in db.
func GetCategories(db *sql.DB, start, count int) ([]Categories, error) {
	rows, err := db.Query(
		"SELECT id, name, created_at FROM categories LIMIT $1 OFFSET $2",
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
		if err := rows.Scan(&dt.ID, &dt.Name, &dt.CreatedAt); err != nil {
			return nil, err
		}
		category = append(category, dt)
	}

	return category, nil
}


// Create new category and insert to database.
func (dt *Categories) CreateCategory(db *sql.DB) error {
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO categories(name, created_at) VALUES($1, $2) RETURNING id, name, created_at", dt.Name, timestamp).Scan(&dt.ID, &dt.Name, &dt.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (cat *Categories) Validate() {
	if cat.Name == "" {
		log.Println("category is required")
	}
}

