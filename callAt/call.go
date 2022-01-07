package callAt

import (
	"fmt"
	"github.com/library/mail"
	"log"
	"time"

)


const timeLayout = "2021-12-22 16:49:31.000000"

func CallAt(callTime string, f func()) error {
	ctime, err := time.Parse(timeLayout, callTime)
	if err != nil {
		return err
	}


	duration := ctime.Sub(time.Now())

	go func() {
		time.Sleep(duration)
		f()
	}()

	return nil
}

func Email() {
	email := mail.NewEmail("", "golang mail", "please, return books to the library")
	err := mail.SendEmail(email)
	log.Print(err)
}

func Task() {
	err := CallAt("2021-12-22 16:49:31.000000", Email)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
