package callAt

import (
	"fmt"
	"github.com/library/db"
	"github.com/library/mail"
	"log"
	"time"
)


var d db.DB

func Email(to []string) {
	email := mail.NewEmail(to, "golang mail", "please, return books to the library")
	err := mail.SendEmail(email)
	log.Print(err)
}

func CheckReturnDate(r db.DB) ([]string, error) {
	transaction, err := r.Database.Begin()
	if err != nil {
		log.Fatalf("Can not begin transaction:%s", err)
		return nil, err
	}
	var listEmail []string
	query := fmt.Sprint("SELECT users.email FROM users JOIN issue ON issue.return_date < $1 AND users.id = issue.user_id")
	rows, err := transaction.Query(query, time.Now().Add(time.Hour))
	if err != nil {
		log.Fatalf("Can not executes a query:%s", err)
		return nil, err
	}
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Fatalf("Scan error:%s", err)
			return nil, err
		}
		listEmail = append(listEmail, email)
	}
	return listEmail, nil
}