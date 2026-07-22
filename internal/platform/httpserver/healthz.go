package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const pingTimeout = 3 * time.Second

// Pinger is the minimal capability healthz needs from a DB pool. *sql.DB
// satisfies it; tests can fake it without a real database.
type Pinger interface {
	PingContext(ctx context.Context) error
}

func healthzHandler(pinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)
		defer cancel()

		status := http.StatusOK
		body := map[string]string{"status": "ok"}
		if err := pinger.PingContext(ctx); err != nil {
			status = http.StatusServiceUnavailable
			body = map[string]string{"status": "unavailable"}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(body)
	}
}
