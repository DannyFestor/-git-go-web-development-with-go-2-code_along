package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string // template.HTML would not be escaped
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Danny Festor",
		Bio:  `<script>alert('HAHAHA')</script>`,
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
