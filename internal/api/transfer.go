package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Represents the payload for a transfer
type TransferRequest struct {
	FromAccountNumber string  `json:"from_account_number"`
	ToAccountNumber   string  `json:"to_account_number"`
	Amount            float64 `json:"amount"`
}

func transferHandler(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := transferFunds(dbpool, userID, req.FromAccountNumber, req.ToAccountNumber, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Success"}`))
}

// Performs the transfer and logs the transaction
func transferFunds(db *pgxpool.Pool, userID int, fromNumber, toNumber string, amount float64) error {
	ctx := context.Background() // shared context for all queries in this transaction

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check source account ownership and balance
	var fromBalance float64
	var fromID int
	err = tx.QueryRow(ctx, "SELECT id, balance FROM accounts WHERE account_number=$1 AND user_id = $2 FOR UPDATE", fromNumber, userID).Scan(&fromID, &fromBalance)
	if err != nil {
		return fmt.Errorf("source account not found or not owned by user")
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// Get destination account
	var toID int
	err = tx.QueryRow(ctx, "SELECT id FROM accounts WHERE account_number=$1 FOR UPDATE", toNumber).Scan(&toID)
	if err != nil {
		return fmt.Errorf("destination account not found")
	}

	// Deduct from source account
	_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return err
	}

	// Credit destination account
	_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return err
	}

	// Log transaction
	_, err = tx.Exec(ctx, `INSERT INTO transactions (user_id, from_account_id, to_account_id, amount) VALUES ($1, $2, $3, $4)`, userID, fromID, toID, amount)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(ctx)
}
