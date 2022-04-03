package callAt

import (
	"github.com/library/db"
	"github.com/library/mail"
	"log"
	"time"
)


func Email(to []string) {
	email := mail.NewEmail(to, "mail", "please, return books to the library")
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
	query := "SELECT users.email FROM users JOIN issue ON issue.return_date < $1 AND users.id = issue.user_id"
	rows, err := transaction.Query(query, time.Now().Add(time.Hour))
	if err != nil {
		log.Fatalf("Can not executes a query:%s", err)
		return nil, err
	}
	for rows.Next() {
		var email string
		if err2 := rows.Scan(&email); err2 != nil {
			log.Fatalf("Scan error:%s", err2)
			return nil, err2
		}
		listEmail = append(listEmail, email)
	}
	return listEmail, nil
}