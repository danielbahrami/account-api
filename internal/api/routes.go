package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Registers all API endpoints
func SetupRoutes(mux *http.ServeMux, dbpool *pgxpool.Pool) {
	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Ok"))
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
		loginHandler(w, r, dbpool)
	})

	// Accounts endpoints
	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authMiddleware(func(w http.ResponseWriter, r *http.Request) {
				listAccountsHandler(w, r, dbpool)
			})(w, r)
		case http.MethodPost:
			authMiddleware(func(w http.ResponseWriter, r *http.Request) {
				createAccountHandler(w, r, dbpool)
			})(w, r)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	// Transfer endpoint
	mux.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			authMiddleware(func(w http.ResponseWriter, r *http.Request) {
				transferHandler(w, r, dbpool)
			})(w, r)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})
}
