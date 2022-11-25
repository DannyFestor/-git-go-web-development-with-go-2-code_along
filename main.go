package main

import (
	"fmt"
	"net/http"
)

// added modd for dynamic reloading
// go install github.com/cortesi/modd/cmd/modd@latest

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // actually not needed as go tries to auto detect the type
	fmt.Fprintf(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:danny@festor.info\">danny@festor.info</a>.")
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path == "/" {
	// 	homeHandler(w, r)
	// } else if r.URL.Path == "/contact" {
	// 	contactHandler(w, r)
	// } else {
	// 	fmt.Fprintf(w, r.URL.Path)
	// }
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		fmt.Fprintln(w, r.URL.Path)
		fmt.Fprintln(w, r.URL.RawPath)
		fmt.Fprintln(w, r.URL.Path)
		fmt.Fprintln(w, r.URL.RawQuery)
		fmt.Fprintln(w, r.URL.RawFragment)
		// TODO page not found error
	}
}

func main() {
	// http.HandleFunc("/", homeHandler)
	http.HandleFunc("/", pathHandler)
	// http.HandleFunc("/contact", contactHandler)
	// http.HandleFunc("/path", pathHandler)
	fmt.Println("Starting the server on :3000 ...")

	http.ListenAndServe(":3000", nil)
}
