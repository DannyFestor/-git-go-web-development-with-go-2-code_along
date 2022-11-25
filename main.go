package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/danakin/web-dev-with-go-2-code_along/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// added modd for dynamic reloading
// go install github.com/cortesi/modd/cmd/modd@latest

func executeTemplate(w http.ResponseWriter, filepath string) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	// r.With(middleware.Logger).Get("/param/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, chi.URLParam(r, "id"))
	// })
	r.NotFound(func(w http.ResponseWriter, r *http.Request) { // not needed but nice to have
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000 ...")
	http.ListenAndServe(":3000", r)
}
