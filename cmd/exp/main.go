package main

import (
	"html/template"
	"os"
)

type User struct {
	Name      string
	Bio       string // template.HTML would not be escaped
	Age       uint32
	Weight    float32
	Hobbies   []string
	Married   bool
	Knowledge map[string]string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name:    "Danny Festor",
		Bio:     `<script>alert('HAHAHA')</script>`,
		Age:     38,
		Weight:  88.5,
		Hobbies: []string{"Programming", "Bicycle", "Eat"},
		Married: true,
		Knowledge: map[string]string{
			"PHP":        "good",
			"Go":         "I'm learning, ok?",
			"Javascript": "cool",
		},
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
