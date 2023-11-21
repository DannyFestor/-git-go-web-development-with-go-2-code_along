package controllers

import (
	"net/http"

	"github.com/danakin/web-dev-with-go-2-code_along/models"
)

type Gallery struct {
	Templates struct {
		New Template
	}

	GalleryService *models.GalleryService
}

func (g Gallery) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}

	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}
