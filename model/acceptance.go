package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines acceptance model.
type Acceptance struct {
	ID               uuid.UUID `json:"id"       sql:"uuid"`
	UserID           string    `json:"userID" validate:"required" sql:"user_id"`
	BookID           string    `json:"bookID" validate:"required" sql:"book_id"`
	BookCondition    string    `json:"bookCondition" validate:"required" sql:"book_condition"`
	Rating           uint      `json:"rating" validate:"required" sql:"rating"`
	FinalCost        float32   `json:"finalCost" validate:"required" sql:"final_cost"`
	Photo            string    `json:"photo" validate:"required" sql:"photo"`
	CreatedAt        time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations

// Gets a specific acceptance by id.
func (dt *Acceptance) GetAcceptance(db *sql.DB) error {
	return db.QueryRow("SELECT user_id, books_id, book_condition, rating, final_cost, photo, created_at, updated_at FROM acceptance WHERE id=$1",
		dt.ID).Scan(&dt.UserID, &dt.BookID, &dt.BookCondition, &dt.Rating, &dt.FinalCost, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets acceptances. Limit count and start position in db.
func GetAcceptances(db *sql.DB, start, count int) ([]Acceptance, error) {
	rows, err := db.Query(
		"SELECT id, user_id, books_id, book_condition, rating, final_cost, photo, created_at, updated_at FROM acceptance LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	acceptance := []Acceptance{}

	// Store query results into acceptance variable if no errors.
	for rows.Next() {
		var dt Acceptance
		if err := rows.Scan(&dt.ID, &dt.UserID, &dt.BookID, &dt.BookCondition, &dt.Rating, &dt.FinalCost, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		acceptance = append(acceptance, dt)
	}

	return acceptance, nil
}

// CRUD operations

// Create new acceptance and insert to database.
func (dt *Acceptance) CreateAcceptance(db *sql.DB) error {
	// Scan db after creation if acceptance exists using new acceptance id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO acceptance(user_id, books_id, book_condition, rating, final_cost, photo, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, user_id, books_id, book_condition, rating, final_cost, photo, created_at, updated_at", dt.UserID, dt.BookID, dt.BookCondition, dt.Rating, dt.FinalCost, dt.Photo, timestamp, timestamp).Scan(&dt.ID, &dt.UserID, &dt.BookID, &dt.BookCondition, &dt.Rating, &dt.FinalCost, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific acceptance details by id.
func (dt *Acceptance) UpdateAcceptance(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE acceptance SET user_id=$1, books_id=$2, book_condition=$3, rating=$4, final_cost=$5, photo=$6, updated_at=$7 WHERE id=$8 RETURNING id, user_id, books_id, book_condition, rating, final_cost, photo, created_at, updated_at",  dt.UserID, dt.BookID, dt.BookCondition, dt.Rating, dt.FinalCost, dt.Photo, timestamp, dt.ID)

	return err
}

// Deletes a specific acceptance by id.
func (dt *Acceptance) DeleteAcceptance(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM acceptance WHERE id=$1", dt.ID)

	return err
}