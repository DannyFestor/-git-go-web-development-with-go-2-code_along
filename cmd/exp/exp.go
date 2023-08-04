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
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err := es.ForgotPassword("danny@festor.info", "https://localhost:3000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent!")
}
