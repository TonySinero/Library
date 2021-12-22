package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines issue model.
type Issue struct {
	ID                uuid.UUID `json:"id"       sql:"uuid"`
	UserID            string    `json:"userid" validate:"required" sql:"userid"`
	BookID            string    `json:"bookid" validate:"required" sql:"bookid"`
	ReturnDate        string    `json:"returndate" validate:"required" sql:"returndate"`
	PreliminaryCost   float32   `json:"preliminarycost" validate:"required" sql:"preliminarycost"`
	CreatedAt         time.Time `json:"createdat" sql:"createdat"`
}

// Query operations

// Gets a specific issue by id.
func (dt *Issue) GetIssue(db *sql.DB) error {
	return db.QueryRow("SELECT userID, bookID, returndate, preliminarycost, createdat FROM issue WHERE id=$1",
		dt.ID).Scan(&dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt)
}

// Gets issuing. Limit count and start position in db.
func GetIssuing(db *sql.DB, start, count int) ([]Issue, error) {
	rows, err := db.Query(
		"SELECT id, userID, bookID, returndate, preliminarycost, createdat FROM issue LIMIT $1 OFFSET $2",
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
		if err := rows.Scan(&dt.ID, &dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt); err != nil {
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
		"INSERT INTO issue(userID, bookID, returndate, preliminarycost, createdat) VALUES($1, $2, $3, $4, $5) RETURNING id, userID, bookID, returndate, preliminarycost, createdat", dt.UserID, dt.BookID, dt.ReturnDate, dt.PreliminaryCost, timestamp).Scan(&dt.ID, &dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt)
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