package controllers

import (
	"fmt"
	"net/http"
)

type User struct {
	Templates struct {
		New Template
	}
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}

func (u User) Store(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Email: ", r.FormValue("email"))
	fmt.Fprintln(w)
	fmt.Fprint(w, "Password: ", r.FormValue("password"))
}
