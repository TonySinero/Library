package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Defines admin model.
type Admin struct {
	ID        uuid.UUID `json:"id" sql:"uuid"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	CreatedAt time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations

// Gets a specific admin by id.
func (u *Admin) GetAdmin(db *sql.DB) error {
	return db.QueryRow("SELECT email, created_at, updated_at FROM admins WHERE id=$1",
		u.ID).Scan(&u.Email, &u.CreatedAt, &u.UpdatedAt)
}

// Gets a specific admin by email and password.
func (u *Admin) GetAdminByEmailAndPassword(db *sql.DB) error {
	return db.QueryRow("SELECT id, email, password, created_at, updated_at FROM admins WHERE email=$1 AND password=$2", u.Email, u.Password).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets multiple admin. Limit count and start position in db.
func GetAdmins(db *sql.DB, field, sort string, limit, page int) ([]Admin, error) {

	rows, err := db.Query( "SELECT id, email, created_at, updated_at FROM admins ORDER BY $1 ,$2 LIMIT $3 OFFSET $4",
		field ,sort ,limit, limit*(page-1))

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	admins := []Admin{}

	// Store query results into admin variable if no errors.
	for rows.Next() {
		var u Admin
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		admins = append(admins, u)
	}

	return admins, nil
}

// CRUD operations

// Create new admin and insert to database.
func (u *Admin) CreateAdmin(db *sql.DB) error {
	if !strings.Contains(u.Email, "@") {
		return errors.New("Email address is required")
	}
	if len(u.Password) < 6 {
		return errors.New("Password is required")
	}
	if len(u.Password) == 0 {
		return errors.New("Password is required")
	}
	temp := &Admin{}
	if temp.Email != "" {
		return errors.New("Email address already in use by another user.")
	}
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO admins(email, password, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id, email, password, created_at, updated_at", u.Email, u.Password, timestamp, timestamp).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific admin's details by id.
func (u *Admin) UpdateAdmin(db *sql.DB) error {
	if !strings.Contains(u.Email, "@") {
		return errors.New("Email address is required")
	}
	if len(u.Password) < 6 {
		return errors.New("Password is required")
	}
	if len(u.Password) == 0 {
		return errors.New("Password is required")
	}
	temp := &Admin{}
	if temp.Email != "" {
		return errors.New("Email address already in use by another user.")
	}
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE admins SET email=$1, password=$2, updated_at=$3 WHERE id=$4 RETURNING id, email, password, created_at, updated_at", u.Email, u.Password, timestamp, u.ID)

	return err
}

// Deletes a specific admin by id.
func (u *Admin) DeleteAdmin(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM admins WHERE id=$1", u.ID)

	return err
}
