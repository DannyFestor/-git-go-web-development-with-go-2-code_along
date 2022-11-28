package controllers

import (
	"net/http"

	"github.com/danakin/web-dev-with-go-2-code_along/views"
)

type User struct {
	Templates struct {
		New views.Template
	}
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
