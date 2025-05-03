package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"titles.run/site_api/models"

	"titles.run/turnstile"
)

func (h *Handler) PostContactHandler() http.HandlerFunc {
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

		var body models.Contact
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := h.core.PostContact(body); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
