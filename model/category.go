package model

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
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
func GetCategories(db *sql.DB, field, sort string, limit, page int) ([]Categories, error) {

	rows, err := db.Query(fmt.Sprintf(  "SELECT id, name, created_at FROM categories ORDER BY %s %s LIMIT %d OFFSET %d",
		field ,sort ,limit, limit*(page-1)))

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

func (dt *Categories) Validate() error {
	if dt.Name == "" {
		return errors.New("category is required")
	}
	return nil
}

