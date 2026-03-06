package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain text password and returns a secure bcrypt hash
func HashPassword(password string) (string, error) {
	// 14 is the "cost" factor. Higher means slower and more secure. 14 is a modern default.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a plain text password with a hashed password from the DB
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT creates a token for a user that expires in 72 hours.
func GenerateJWT(userID uuid.UUID) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET environment variable is not set.")
	}

	// Create the token claims (the data stored inside the token)
	claims := jwt.MapClaims{
		"userID": userID.String(),
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Expires in 3 days
		"iat":    jwt.NewNumericDate(time.Now()),                     // Issued at
	}

	// Create the token with the claims and the signing method (HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString([]byte(secretKey))
}

// ValidateJWT parses a token string and return the UserID if valid
func ValidateJWT(tokenString string) (uuid.UUID, error) {
	secretKey := os.Getenv("JWT_SECRET")

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("enexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("token parsed successfully but is marked invalid")
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	// Extract the user ID from the claims
	userIDString, ok := claims["userID"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in token")
	}

	// Convertstring back to UUID
	parsedUUID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, errors.New("invalid user_id format in token")
	}

	return parsedUUID, nil
}
