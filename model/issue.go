package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines issue model.
type Issue struct {
	ID                uuid.UUID `json:"id"       sql:"uuid"`
	UserID            string    `json:"userID" validate:"required" sql:"user_id"`
	BooksID           string    `json:"bookID" validate:"required" sql:"book_id"`
	ReturnDate        string    `json:"returnDate" validate:"required" sql:"return_date"`
	PreliminaryCost   float32   `json:"preliminaryCost" validate:"required" sql:"preliminary_cost"`
	CreatedAt         time.Time `json:"createdAt" sql:"created_at"`
}

// Query operations

// Gets a specific issue by id.
func (dt *Issue) GetIssue(db *sql.DB) error {
	return db.QueryRow("SELECT user_id, books_id, return_date, preliminary_cost, created_at FROM issue WHERE id=$1",
		dt.ID).Scan(&dt.UserID, &dt.BooksID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt)
}

// Gets issuing. Limit count and start position in db.
func GetIssuing(db *sql.DB, start, count int) ([]Issue, error) {
	rows, err := db.Query(
		"SELECT id, user_id, books_id, return_date, preliminary_cost, created_at FROM issue LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	issue := []Issue{}

	// Store query results into lending variable if no errors.
	for rows.Next() {
		var dt Issue
		if err := rows.Scan(&dt.ID, &dt.UserID, &dt.BooksID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt); err != nil {
			return nil, err
		}
		issue = append(issue, dt)
	}

	return issue, nil
}

// CRUD operations

// Create new issue and insert to database.
func (dt *Issue) CreateIssue(db *sql.DB) error {
	// Scan db after creation if issue exists using new issue id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO issue(user_id, books_id, return_date, preliminary_cost, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id, user_id, books_id, return_date, preliminary_cost, created_at", dt.UserID, dt.BooksID, dt.ReturnDate, dt.PreliminaryCost, timestamp).Scan(&dt.ID, &dt.UserID, &dt.BooksID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}


// Deletes a specific issue by id.
func (dt *Issue) DeleteIssue(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM issue WHERE id=$1", dt.ID)

	return err
}