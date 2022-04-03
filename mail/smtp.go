package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

const (
	HOST        = "smtp.gmail.com"
	PORT        = "587"
	USER        = ""
	PASSWORD    = ""
)

type Email struct {
	to      []string "to"
	subject   string "subject"
	msg       string "msg"
}

func NewEmail(to []string, subject, msg string) *Email {
	return &Email{to:to, subject: subject, msg: msg}
}
func SendEmail(email *Email) error {
	auth := smtp.PlainAuth("", USER, PASSWORD, HOST)
	sendTo := email.to
	addr := fmt.Sprintf("%s:%s", HOST, PORT)

		for _, v := range sendTo {
			str := strings.Replace("From: "+USER+"~To: "+v+"~Subject: "+email.subject+"~~", "~", "\r\n", -1) + email.msg
			err := smtp.SendMail(addr, auth, USER,	[]string{v}, []byte(str))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Successfully sent mail to all user in toList")
		}
	return nil
}

