package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines author model.
type Author struct {
	ID             uuid.UUID `json:"id"       sql:"uuid"`
	Firstname      string    `json:"firstname" validate:"required" sql:"firstname"`
	Surname        string    `json:"surname" validate:"required" sql:"surname"`
	DateOfBirth    string    `json:"dateofbirth" validate:"required" sql:"dateofbirth"`
	Photo          string    `json:"photo" validate:"required" sql:"photo"`
	CreatedAt      time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt      time.Time `json:"updatedat" sql:"updatedat"`
}

// Query operations

// Gets authors. Limit count and start position in db.
func GetAuthors(db *sql.DB, start, count int) ([]Author, error) {
	rows, err := db.Query(
		"SELECT id, firstname, surname, dateofbirth, photo, createdat, updatedat FROM authors LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	author := []Author{}

	// Store query results into author variable if no errors.
	for rows.Next() {
		var dt Author
		if err := rows.Scan(&dt.ID, &dt.Firstname, &dt.Surname, &dt.DateOfBirth, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		author = append(author, dt)
	}

	return author, nil
}

// CRUD operations

// Create new author and insert to database.
func (dt *Author) CreateAuthor(db *sql.DB) error {
	// Scan db after creation if author exists using new author id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO authors(firstname, surname, dateofbirth, photo, createdat, updatedat) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, firstname, surname, dateofbirth, photo, createdat, updatedat", dt.Firstname, dt.Surname,  dt.DateOfBirth, dt.Photo, timestamp, timestamp).Scan(&dt.ID, &dt.Firstname, &dt.Surname, &dt.DateOfBirth, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific author details by id.
func (dt *Author) UpdateAuthor(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE authors SET firstname=$1, surname=$2, dateofbirth=$3, photo=$4, updatedat=$5 WHERE id=$6 RETURNING id, firstname, surname, photo, createdat, updatedat", dt.Firstname, dt.Surname,  dt.DateOfBirth, dt.Photo, timestamp, dt.ID)

	return err
}

