package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetWebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()

		if err := h.titles.GetWebhook(queryParams.Get("hub.verify_token")); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]string{
				"hub.challenge": queryParams.Get("hubChallenge"),
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		http.Error(w, "Invalid verification token", http.StatusForbidden)
	}
}
