package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"titles.run/turnstile"
	"titles.run/webhook/models"
)

func (h *Handler) PostDemoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Titles-CF-Turnstile-Token")

		turnstile := turnstile.TurnStile{
			Secret: os.Getenv("CF_TURNSTILE_SECRETKEY"),
			Token:  token,
		}

		err := turnstile.Verify()
		if err != nil {
			http.Error(w, "Error validating turnstile token", http.StatusBadRequest)
			return
		}

		var body models.Gpx
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		title, err := h.titles.PostDemo(body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(title))
	}
}
