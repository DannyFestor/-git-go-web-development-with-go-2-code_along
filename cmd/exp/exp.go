package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/danakin/web-dev-with-go-2-code_along/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("MAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		panic(err)
	}
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err = es.ForgotPassword("danny@festor.info", "https://localhost:3000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent!")
}
