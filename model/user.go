package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Defines user model.
type User struct {
	ID             uuid.UUID `json:"id"       sql:"uuid"`
	Firstname      string    `json:"firstName" validate:"required" sql:"firstname"`
	Surname        string    `json:"surname" validate:"required" sql:"surname"`
	SecondName     string    `json:"secondName" validate:"required" sql:"second_name"`
	Passport       string    `json:"passport" validate:"required" sql:"passport"`
	DateOfBirth    string    `json:"dateOfBirth" validate:"required" sql:"date_of_birth"`
	Email          string    `json:"email" validate:"required" sql:"email"`
	Address        string    `json:"address" validate:"required" sql:"address"`
	Indebtedness   string    `json:"indebtedness" validate:"required" sql:"indebtedness"`
	CreatedAt      time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations

// Gets a specific user by id.
func (dt *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT firstname, surname, second_name, passport, date_of_birth, email, address, indebtedness, created_at, updated_at FROM users WHERE id=$1",
		dt.ID).Scan(&dt.Firstname, &dt.Surname, &dt.SecondName, &dt.Passport, &dt.DateOfBirth, &dt.Email, &dt.Address, &dt.Indebtedness, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets users. Limit count and start position in db.
func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, firstname, surname, second_name, passport, date_of_birth, email, address, indebtedness, created_at, updated_at FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	user := []User{}

	// Store query results into user variable if no errors.
	for rows.Next() {
		var dt User
		if err := rows.Scan(&dt.ID, &dt.Firstname, &dt.Surname, &dt.SecondName, &dt.Passport, &dt.DateOfBirth, &dt.Email, &dt.Address, &dt.Indebtedness, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		user = append(user, dt)
	}

	return user, nil
}

// CRUD operations

// Create new user and insert to database.
func (dt *User) CreateUser(db *sql.DB) error {
	// Scan db after creation if user exists using new user id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO users(firstname, surname, second_name, passport, date_of_birth, email, address, indebtedness, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, firstname, surname, second_name, passport, date_of_birth, email, address, indebtedness, created_at, updated_at", dt.Firstname, dt.Surname, dt.SecondName, dt.Passport, dt.DateOfBirth, dt.Email, dt.Address, dt.Indebtedness, timestamp, timestamp).Scan(&dt.ID, &dt.Firstname, &dt.Surname, &dt.SecondName, &dt.Passport, &dt.DateOfBirth, &dt.Email, &dt.Address, &dt.Indebtedness, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific user details by id.
func (dt *User) UpdateUser(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE users SET firstname=$1, surname=$2, second_name=$3, passport=$4, date_of_birth=$5, email=$6, address=$7, indebtedness=$8, updated_at=$9 WHERE id=$10 RETURNING id, firstname, surname, second_name, passport, date_of_birth, email, address, indebtedness, created_at, updated_at", dt.Firstname, dt.Surname, dt.SecondName, dt.Passport, dt.DateOfBirth, dt.Email, dt.Address, dt.Indebtedness, timestamp, dt.ID)

	return err
}

// Deletes a specific user by id.
func (dt *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", dt.ID)

	return err
}

func (user *User) Validate() error {
	if user.Firstname == "" {
		return errors.New("name is required")
	}
	if user.Surname == "" {
		return errors.New("surname is required")
	}
	if user.SecondName == "" {
		return errors.New("secondName is required")
	}
	if user.DateOfBirth == "" {
		return errors.New("dateOfBirth is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("Email address is required")
	}
	if user.Address == "" {
		return errors.New("address is required")
	}
	if user.Indebtedness == "" {
		return errors.New("indebtedness is required")
	}
	return nil
}