package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
)


// Defines issue model.

type Issue struct {
	ID                uuid.UUID `json:"id"       sql:"uuid"`
	UserID            uuid.UUID `json:"userID" validate:"required" sql:"user_id"`
	BookID            uuid.UUID `json:"bookID" validate:"required" sql:"book_id"`
	ReturnDate        string    `json:"returnDate" validate:"required" sql:"return_date"`
	PreliminaryCost   float32   `json:"preliminaryCost" validate:"required" sql:"preliminary_cost"`
	CreatedAt         time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations

// Gets a specific user by id.
func (dt *Issue) GetIssue(db *sql.DB) error {
	return db.QueryRow("SELECT user_id, book_id, return_date, preliminary_cost, created_at, updated_at FROM issue WHERE id=$1",
		dt.ID).Scan(&dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets users. Limit count and start position in db.
func GetIssues(db *sql.DB, field, sort string, limit, page int) ([]Issue, error) {

	rows, err := db.Query(  "SELECT * FROM issue ORDER BY $1 ,$2 LIMIT $3 OFFSET $4",
		field ,sort ,limit, limit*(page-1))

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	issue := []Issue{}

	// Store query results into user variable if no errors.
	for rows.Next() {
		var dt Issue
		if err := rows.Scan(&dt.ID, &dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		issue = append(issue, dt)
	}

	return issue, nil
}

// CRUD operations

// Create new user and insert to database.
func (dt *Issue) CreateIssue(db *sql.DB) error {
	if dt.ReturnDate == "" {
		return errors.New("date is required")
	}
	if dt.PreliminaryCost == 0 {
		return errors.New("cost cannot be zero")
	}

	// Scan db after creation if user exists using new user id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO issue(user_id, book_id, return_date, preliminary_cost, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, user_id, book_id, return_date, preliminary_cost, created_at, updated_at", &dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, timestamp, timestamp).Scan(&dt.ID, &dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Updates a specific user details by id.
func (dt *Issue) UpdateIssue(db *sql.DB) error {
	if dt.ReturnDate == "" {
		return errors.New("date is required")
	}
	if dt.PreliminaryCost == 0 {
		return errors.New("cost cannot be zero")
	}
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE issue SET user_id=$1, book_id=$2, return_date=$3, preliminary_cost=$4, updated_at=$5 WHERE id=$6 RETURNING id, user_id, book_id, return_date, preliminary_cost, created_at, updated_at", &dt.UserID, &dt.BookID, &dt.ReturnDate, &dt.PreliminaryCost, timestamp, dt.ID)

	return err
}

// Deletes a specific user by id.
func (dt *Issue) DeleteIssue(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM issue WHERE id=$1", dt.ID)

	return err
}

func (dt *Issue) PremCostFunc(b *Book) {
	dt.PreliminaryCost = b.PricePerDay * 30
}


