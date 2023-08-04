package main

import (
	"fmt"
	"os"

	"github.com/go-mail/mail/v2"
)

const (
	host     = "localhost" // sandbox.smtp.mailtrap.io etc.
	port     = 1025        // standard port is 587
	username = "mailpit"   // mailtrap username
	password = "mailpit"   // mailtrap password
)

func main() {
	from := "test@lenslocked.com"
	to := "danny@festor.info"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>I hope it finds you well!</p>`

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)

	msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(host, port, username, password)
	// Approach 1: create a sender and a dialer, to send multiple emails from the same sender
	sender, err := dialer.Dial()
	if err != nil {
		panic(err)
	}
	defer sender.Close()
	sender.Send(from, []string{to}, msg)

	// Approach 2: dial and send directly
	err = dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent")
}
