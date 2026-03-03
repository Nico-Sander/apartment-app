// cmd/api/main.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"apartment-app/db"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// 2. Connect to Database
	db.InitDB()
	defer db.Pool.Close()

	// 3. Serve Static Files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 4. Parse Templates
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))

	// 5. Routes

	// Homepage Route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "base.html", nil)
	})

	// Create User Route (Receives the HTMX POST request)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data sent by the browser
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Extract the values
		name := r.FormValue("name")
		email := r.FormValue("email")

		// Call our database function!
		user, err := db.CreateUser(name, email)
		if err != nil {
			// If the email already exists, Postgres will throw an error
			errorHtml := fmt.Sprintf(`<p class="text-red-500 font-medium">Error: Could not create user. Does the email %s already exist?</p>`, email)
			w.Write([]byte(errorHtml))
			return
		}

		// If successful, return a success message in HTML format
		successHtml := fmt.Sprintf(`
			<div class="p-4 bg-green-50 border border-green-200 rounded text-green-800">
				<p class="font-bold">✅ User Created Successfully!</p>
				<p class="text-sm mt-1">Name: %s</p>
				<p class="text-sm">ID: <span class="font-mono text-xs">%s</span></p>
			</div>
		`, user.Name, user.ID.String())

		w.Write([]byte(successHtml))
	})

	// 6. Start Server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
