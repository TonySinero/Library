package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
)

// Defines popular model.
type Popular struct {
	ID               uuid.UUID `json:"id"       sql:"uuid"`
	BookID           uuid.UUID `json:"bookID" validate:"required" sql:"book_id"`
	Rating           uint      `json:"rating" validate:"required" sql:"rating"`
	Views            uint      `json:"views"  validate:"required" sql:"views"`
	CreatedAt        time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations


// Gets popular. Limit count and start position in db.
func GetPopulars(db *sql.DB) ([]Popular, error) {
	rows, err := db.Query(
		"SELECT id, book_id, rating, views, created_at, updated_at FROM popular ORDER BY views DESC LIMIT 3",
		)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	popular := []Popular{}

	// Store query results into user variable if no errors.
	for rows.Next() {
		var dt Popular
		if err := rows.Scan(&dt.ID, &dt.BookID, &dt.Rating, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		popular = append(popular, dt)
	}
	return popular, nil
}

// CRUD operations

// Create new popular and insert to database.
func (dt *Popular) CreatePopular(db *sql.DB) error {
	// Scan db after creation if popular exists using new popular id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO popular(book_id, rating, views, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id, book_id, rating, views, created_at, updated_at", dt.BookID, dt.Rating, dt.Views, timestamp, timestamp).Scan(&dt.ID, &dt.BookID, &dt.Rating, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Updates a specific popular details by id.
func (dt *Popular) UpdatePopular(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE popular SET book_id=$1, rating=$2, views=$3, updated_at=$4 WHERE id=$5 RETURNING id, book_id, rating, views, created_at, updated_at",  dt.BookID, dt.Rating, dt.Views, timestamp, dt.ID)

	return err
}

// Deletes a specific popular by id.
func (dt *Popular) DeletePopular(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM popular WHERE id=$1", dt.ID)

	return err
}

func (dt *Popular) Validate() error {
	if dt.Rating == 0 {
		return errors.New("rating cannot be zero")
	}
	if dt.Views == 0 {
		return errors.New("photo is required")
	}
	return nil
}
