package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Represents the payload from the client
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Represents the response sent to the client
type LoginResponse struct {
	Token string `json:"token"`
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Handles user authentication and JWT generation
func loginHandler(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	var id int
	var passwordHash string

	sql := "SELECT id, password_hash FROM users WHERE email = $1"

	row := dbpool.QueryRow(context.Background(), sql, req.Email)

	err := row.Scan(&id, &passwordHash)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(time.Minute * 10).Unix(), // Token expires after 10 minutes
		})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
