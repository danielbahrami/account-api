package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

}
