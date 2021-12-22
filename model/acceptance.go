package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines acceptance model.
type Acceptance struct {
	ID               uuid.UUID `json:"id"       sql:"uuid"`
	UserID           string    `json:"userid" validate:"required" sql:"userid"`
	BookID           string    `json:"bookid" validate:"required" sql:"bookid"`
	BookCondition    string    `json:"bookcondition" validate:"required" sql:"bookcondition"`
	Rating           uint      `json:"rating" validate:"required" sql:"rating"`
	FinalCost        float32   `json:"finalcost" validate:"required" sql:"finalcost"`
	Photo            string    `json:"photo" validate:"required" sql:"photo"`
	CreatedAt        time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt        time.Time `json:"updatedat" sql:"updatedat"`
}

// Query operations

// Gets a specific acceptance by id.
func (dt *Acceptance) GetAcceptance(db *sql.DB) error {
	return db.QueryRow("SELECT userID, bookID, bookcondition, rating, finalcost, photo, createdat, updatedat FROM acceptance WHERE id=$1",
		dt.ID).Scan(&dt.UserID, &dt.BookID, &dt.BookCondition, &dt.Rating, &dt.FinalCost, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets acceptances. Limit count and start position in db.
func GetAcceptances(db *sql.DB, start, count int) ([]Acceptance, error) {
	rows, err := db.Query(
		"SELECT id, userID, bookID, bookcondition, rating, finalcost, photo, createdat, updatedat FROM acceptance LIMIT $1 OFFSET $2",
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
		"INSERT INTO acceptance(userID, bookID, bookcondition, rating, finalcost, photo, createdat, updatedat) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, userID, bookID, bookcondition, rating, finalcost, photo, createdat, updatedat", dt.UserID, dt.BookID, dt.BookCondition, dt.Rating, dt.FinalCost, dt.Photo, timestamp, timestamp).Scan(&dt.ID, &dt.UserID, &dt.BookID, &dt.BookCondition, &dt.Rating, &dt.FinalCost, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific acceptance details by id.
func (dt *Acceptance) UpdateAcceptance(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE acceptance SET userID=$1, bookID=$2, bookcondition=$3, rating=$4, finalcost=$5, photo=$6, updatedat=$7 WHERE id=$8 RETURNING id, userID, bookID, bookcondition, rating, finalcost, photo, createdat, updatedat",  dt.UserID, dt.BookID, dt.BookCondition, dt.Rating, dt.FinalCost, dt.Photo, timestamp, dt.ID)

	return err
}

// Deletes a specific acceptance by id.
func (dt *Acceptance) DeleteAcceptance(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM acceptance WHERE id=$1", dt.ID)

	return err
}