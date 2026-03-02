// cmd/api/main.go
package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	// 1. Serve static files (like your compiled CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Parse the HTML templates
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))

	// 3. Route: Homepage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "base.html", nil)
	})

	// 4. Route: HTMX endpoint for the button
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong! The Go server is awake and HTMX is working!"))
	})

	// 5. Start the server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}