package main

import (
	"fmt"

	"github.com/danakin/web-dev-with-go-2-code_along/models"
)

const (
	host     = "localhost" // sandbox.smtp.mailtrap.io etc.
	port     = 1025        // standard port is 587
	username = "mailpit"   // mailtrap username
	password = "mailpit"   // mailtrap password
)

func main() {
	email := models.Email{
		From:      "test@lenslocked.com",
		To:        "danny@festor.info",
		Subject:   "This is a test email",
		Plaintext: "This is the body of the email",
		HTML:      `<h1>Hello there buddy!</h1><p>This is the email</p><p>I hope it finds you well!</p>`,
	}
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err := es.Send(email)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent!")
}
