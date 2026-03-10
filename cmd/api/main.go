// cmd/api/main.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"apartment-app/auth"
	"apartment-app/db"

	"github.com/google/uuid"
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

	// 4. Parse Templates into seperate variables
	indexTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	dashboardTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/dashboard.html"))

	// 5. Routes

	// Homepage Route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. Check if the user already has a session cookie (logged in)
		cookie, err := r.Cookie("apartment_session")
		if err == nil {
			// 2. They have a cookie, validate JWT to ensure it's not fake or expired
			_, err = auth.ValidateJWT(cookie.Value)
			if err == nil {
				// 3. Valid user -> Redirect them instantly to the dashboard
				http.Redirect(w, r, "/dashboard", http.StatusFound)
				return
			}
		}

		// 4. If there is no cookie, or it's invalid, show the login, registration page
		err = indexTmpl.ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
			log.Println("Template Error:", err)
		}
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
		password := r.FormValue("password")

		// Call our database function!
		_, token, err := db.CreateUser(name, email, password)
		if err != nil {
			// If the email already exists, Postgres will throw an error
			errorHtml := fmt.Sprintf(`<p class="text-red-500 font-medium">Error: Could not create user. Does the email %s already exist?</p>`, email)
			w.Write([]byte(errorHtml))
			return
		}

		// Set the secure cookie directly upon registration!
		http.SetCookie(w, &http.Cookie{
			Name:     "apartment_session",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   72 * 3600,
		})

		// Redirect the brand new user straight to the dashboard
		w.Header().Set("HX-Redirect", "/dashboard")
		w.Write([]byte("Redirecting to dashboard..."))
	})

	// Login Route
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		// 1. Find hte user in the DB
		user, err := db.GetUserByEmail(email)
		if err != nil {
			w.Write([]byte(`<p class="text-red-500 font-medium">Invalid email or password.</p>`))
			return
		}

		// 2. Verify password
		if !auth.CheckPasswordHash(password, user.PasswordHash) {
			w.Write([]byte(`<p class="text-red-500 font-medium">Invalid email or password.</p>`))
			return
		}

		// 3. Generate JWT
		token, err := auth.GenerateJWT(user.ID)
		if err != nil {
			w.Write([]byte(`<p class="text-red-500 font-medium">Internal server error.</p>`))
			return
		}

		// 4. Set the Secure Cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "apartment_session",
			Value:    token,
			Path:     "/",                  // Cookie is valid for the whole website
			HttpOnly: true,                 // JavaScript cannot access this
			Secure:   false,                // TODO: Set 'true' in production when using HTTPS!
			SameSite: http.SameSiteLaxMode, // Protects against Cross-Site Request Forgery
			MaxAge:   72 * 3600,            // Expires in 72 hours
		})

		// 5. Redirect the user to the dashboard HTMX
		w.Header().Set("HX-Redirect", "/dashboard")
		w.Write([]byte("Redirecting to dashboard..."))
	})

	// Logout Route
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// Overwrite the existing cookie with a dead one
		http.SetCookie(w, &http.Cookie{
			Name:     "apartment_session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1, // This tells the browser to delete the cookie immediately
		})

		// Redirect back to the home page
		w.Header().Set("HX-Redirect", "/")

		// Fallback  message just in case
		w.Write([]byte("Logged out successfully."))
	})

	// Protected Dashboard Route
	http.HandleFunc("/dashboard", auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		// Extract the securely validated UserID from the request context
		userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

		// Execute the dashboard template and pass the userID into the HTML
		err := dashboardTmpl.ExecuteTemplate(w, "base.html", userID.String())
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
		}
	}))

	// 6. Start Server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
