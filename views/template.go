package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/danakin/web-dev-with-go-2-code_along/context"
	"github.com/danakin/web-dev-with-go-2-code_along/models"
	"github.com/gorilla/csrf"
)

type public interface {
	Public() string // errors that have a Public() method will be passed to the views
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// templates are parsed from base name only, so having a/foo.gohtml and b/foo.gohtml will be parsed as foo and only b/foo.gohtml will be available; https://pkg.go.dev/html/template#ParseFiles
	tpl := template.New(path.Base(patterns[0]))
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) { // gets overwritten in execute because we don't have access to request at this point (this gets run at the start of the app instead of request so we don't have to parse all templates at every request)
				return "", fmt.Errorf("csrfField template function not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser template function not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing FS template: %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, err
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, err
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone() // clone the template to ensure each users gets a different one when users hit the page at the same time
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}

	errorMessages := errMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML { // update csrfField stub to actually use TemplateField
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errorMessages
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// add function because closure inside template functions are run twice atm...
func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) { // only add errors that implement the public interface to the views
			msgs = append(msgs, pubErr.Public())
		} else { // otherwise add a generic error
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
