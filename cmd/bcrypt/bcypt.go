package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Invalid command: %v\n", os.Args[1])
	}
}

func hash(password string) {
	fmt.Printf("Hashing the password %q\n", password)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing password")
		return
	}
	fmt.Printf("Hashed password %q\n", hashedBytes)
}

func compare(password string, hash string) {
	fmt.Printf("TODO: Compare the password %q with the hash %q\n", password, hash)
}
