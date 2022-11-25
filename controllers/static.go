package controllers

import (
	"net/http"

	"github.com/danakin/web-dev-with-go-2-code_along/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}
