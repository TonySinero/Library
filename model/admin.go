package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines admin model.
type Admin struct {
	ID        uuid.UUID `json:"id" sql:"uuid"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	CreatedAt time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt time.Time `json:"updatedat" sql:"updatedat"`
}

// Query operations

// Gets a specific admin by id.
func (u *Admin) GetAdmin(db *sql.DB) error {
	return db.QueryRow("SELECT email, password, createdat, updatedat FROM admins WHERE id=$1",
		u.ID).Scan(&u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets a specific admin by email and password.
func (u *Admin) GetAdminByEmailAndPassword(db *sql.DB) error {
	return db.QueryRow("SELECT id, email, password, createdat, updatedat FROM admins WHERE email=$1 AND password=$2", u.Email, u.Password).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets multiple admin. Limit count and start position in db.
func GetAdmins(db *sql.DB, start, count int) ([]Admin, error) {
	rows, err := db.Query(
		"SELECT id, email, password, createdat, updatedat FROM admins LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	admins := []Admin{}

	// Store query results into admin variable if no errors.
	for rows.Next() {
		var u Admin
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		admins = append(admins, u)
	}

	return admins, nil
}

// CRUD operations

// Create new admin and insert to database.
func (u *Admin) CreateAdmin(db *sql.DB) error {
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO admins(email, password, createdat, updatedat) VALUES($1, $2, $3, $4) RETURNING id, email, password, createdat, updatedat", u.Email, u.Password, timestamp, timestamp).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific admin's details by id.
func (u *Admin) UpdateAdmin(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE admins SET email=$1, password=$2, updatedat=$3 WHERE id=$4 RETURNING id, email, password, createdat, updatedat", u.Email, u.Password, timestamp, u.ID)

	return err
}

// Deletes a specific admin by id.
func (u *Admin) DeleteAdmin(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM admins WHERE id=$1", u.ID)

	return err
}
