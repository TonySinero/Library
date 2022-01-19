package model

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
)

// Defines author model.
type Author struct {
	ID             uuid.UUID `json:"id"       sql:"uuid"`
	Firstname      string    `json:"firstname" validate:"required" sql:"firstname"`
	Surname        string    `json:"surname" validate:"required" sql:"surname"`
	DateOfBirth    string    `json:"dateOfBirth" validate:"required" sql:"date_of_birth"`
	Photo          string    `json:"photo" validate:"required" sql:"photo"`
	CreatedAt      time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations

// Gets authors. Limit count and start position in db.
func GetAuthors(db *sql.DB, field, sort string, limit, page int) ([]Author, error) {

	rows, err := db.Query(fmt.Sprintf(  "SELECT id, firstname, surname, date_of_birth, photo, created_at, updated_at FROM authors ORDER BY %s %s LIMIT %d OFFSET %d",
		field ,sort ,limit, limit*(page-1)))

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
		"INSERT INTO authors(firstname, surname, date_of_birth, photo, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, firstname, surname, date_of_birth, photo, created_at, updated_at", dt.Firstname, dt.Surname,  dt.DateOfBirth, dt.Photo, timestamp, timestamp).Scan(&dt.ID, &dt.Firstname, &dt.Surname, &dt.DateOfBirth, &dt.Photo, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific author details by id.
func (dt *Author) UpdateAuthor(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE authors SET firstname=$1, surname=$2, date_of_birth=$3, photo=$4, updated_at=$5 WHERE id=$6 RETURNING id, firstname, surname, date_of_birth, photo, created_at, updated_at", dt.Firstname, dt.Surname,  dt.DateOfBirth, dt.Photo, timestamp, dt.ID)

	return err
}

func (dt *Author) Validate() error {
	if dt.Firstname == "" {
		return errors.New("name is required")
	}
	if dt.Surname == "" {
		return errors.New("surname is required")
	}
	if dt.DateOfBirth == "" {
		return errors.New("dateOfBirth is required")
	}
	if dt.Photo == "" {
		return errors.New("photo is required")
	}
	return nil
}