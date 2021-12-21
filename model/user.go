package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines user model.
type User struct {
	ID             uuid.UUID `json:"id"       sql:"uuid"`
	Firstname      string    `json:"firstname" validate:"required" sql:"firstname"`
	Surname        string    `json:"surname" validate:"required" sql:"surname"`
	SecondName     string    `json:"secondname" validate:"required" sql:"secondname"`
	Passport       string    `json:"passport" validate:"required" sql:"passport"`
	DateOfBirth    string    `json:"dateofbirth" validate:"required" sql:"dateofbirth"`
	Email          string    `json:"email" validate:"required" sql:"email"`
	Address        string    `json:"address" validate:"required" sql:"address"`
	Indebtedness   string    `json:"indebtedness" validate:"required" sql:"indebtedness"`
	CreatedAt      time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt      time.Time `json:"updatedat" sql:"updatedat"`
}

// Query operations

// Gets a specific user by id.
func (dt *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT firstname, surname, secondname, passport, dateofbirth, email, address, indebtedness, createdat, updatedat FROM users WHERE id=$1",
		dt.ID).Scan(&dt.Firstname, &dt.Surname, &dt.SecondName, &dt.Passport, &dt.DateOfBirth, &dt.Email, &dt.Address, &dt.Indebtedness, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets users. Limit count and start position in db.
func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, firstname, surname, secondname, passport, dateofbirth, email, address, indebtedness, createdat, updatedat FROM users LIMIT $1 OFFSET $2",
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
		"INSERT INTO users(firstname, surname, secondname, passport, dateofbirth, email, address, indebtedness, createdat, updatedat) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, firstname, surname, secondname, passport, dateofbirth, email, address, indebtedness, createdat, updatedat", dt.Firstname, dt.Surname, dt.SecondName, dt.Passport, dt.DateOfBirth, dt.Email, dt.Address, dt.Indebtedness, timestamp, timestamp).Scan(&dt.ID, &dt.Firstname, &dt.Surname, &dt.SecondName, &dt.Passport, &dt.DateOfBirth, &dt.Email, &dt.Address, &dt.Indebtedness, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific user details by id.
func (dt *User) UpdateUser(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE users SET firstname=$1, surname=$2, secondname=$3, passport=$4, dateofbirth=$5, email=$6, address=$7, indebtedness=$8, updatedat=$9 WHERE id=$10 RETURNING id, firstname, surname, secondname, passport, dateofbirth, email, address, indebtedness, createdat, updatedat", dt.Firstname, dt.Surname, dt.SecondName, dt.Passport, dt.DateOfBirth, dt.Email, dt.Address, dt.Indebtedness, timestamp, dt.ID)

	return err
}

// Deletes a specific user by id.
func (dt *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", dt.ID)

	return err
}
