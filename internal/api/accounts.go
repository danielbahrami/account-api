package api

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Response payload to client
type AccountResponse struct {
	AccountName   string  `json:"account_name"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

// Request payload from client
type CreateAccountRequest struct {
	AccountName string `json:"account_name"`
}

func listAccountsHandler(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sql := "SELECT account_number, account_name, balance FROM accounts WHERE user_id = $1 ORDER BY account_name ASC"

	rows, err := dbpool.Query(context.Background(), sql, userID)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a slice of AccountResponse
	accounts := []AccountResponse{}
	for rows.Next() {
		var account AccountResponse
		err := rows.Scan(&account.AccountNumber, &account.AccountName, &account.Balance)
		if err != nil {
			http.Error(w, "Failed to parse account data", http.StatusInternalServerError)
			return
		}
		accounts = append(accounts, account)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	// Decode request JSON
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user_id from context
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create account
	acc, err := createAccount(dbpool, userID, req.AccountName)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func createAccount(db *pgxpool.Pool, userID int, accountName string) (AccountResponse, error) {
	for {
		accountNumber, err := generateAccountNumber()
		if err != nil {
			return AccountResponse{}, err
		}

		var acc AccountResponse
		sql := "INSERT INTO accounts (user_id, account_name, account_number, balance) VALUES ($1, $2, $3, 0.00) RETURNING account_name, account_number, balance"
		err = db.QueryRow(context.Background(), sql, userID, accountName, accountNumber).Scan(&acc.AccountName, &acc.AccountNumber, &acc.Balance)

		if err == nil {
			return acc, nil
		}

		// Retry if account_number UNIQUE constraint is violated
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			continue
		}
		return AccountResponse{}, err
	}
}

// Returns a random 10 digit string
func generateAccountNumber() (string, error) {
	var number string
	for range 10 {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		number += fmt.Sprintf("%d", n.Int64())
	}
	return number, nil
}
