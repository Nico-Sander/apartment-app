// cmd/api/main.go
package main

import (
	"html/template"
	"log"
	"net/http"

	"apartment-app/db"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// 2. Connect to the database
	db.InitDB()
	defer db.Pool.Close()

	// 3. Serve static files (like your compiled CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 4. Parse the HTML templates
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))

	// 5. Route: Homepage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "base.html", nil)
	})

	// 6. Route: HTMX endpoint for the button
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong! The Go server is awake and HTMX is working!"))
	})

	// 7. Start the server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

