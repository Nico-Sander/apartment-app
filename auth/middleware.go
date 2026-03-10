package auth

import (
	"context"
	"net/http"
	//"github.com/google/uuid"
)

// Define a custom type for the Context Key.
// Prevents naming collisions if other packages also try to save data in the request context.
type contextKey string

const UserIDKey contextKey = "user_id"

// RequireAuth is the bouncer. It wraps around standard http.HandleFunc routes.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Try to grab the secure cookie from the incoming request
		cookie, err := r.Cookie("apartment_session")
		if err != nil {
			// No cookie -> Bounce
			http.Error(w, "Unauthorized: Please log in first.", http.StatusUnauthorized)
			return
		}

		// 2. Validate the JWT inside the cookie using the existing function
		userID, err := ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid or expired session.", http.StatusUnauthorized)
			return
		}

		// 3. Attach the UserID to the request Context
		// It allows the final route to know exacatly who is making the request without having to parse the token all over again
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)

		// 4. The user is legit. Let them throhg to the actual route
		next.ServeHTTP(w, r)
	}
}
