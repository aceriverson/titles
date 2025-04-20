package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"titles.run/webhook/models"
)

func (h *Handler) PostWebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body models.Webhook
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if strconv.FormatInt(body.SubscriptionID, 10) != os.Getenv("STRAVA_WEBHOOK_SUBSCRIPTION") {
			http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
			return
		}

		if body.Updates["authorized"] == "false" {
			if err := h.titles.UnauthorizeUser(body.OwnerID); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		if body.ObjectType != "activity" || body.AspectType != "create" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := h.titles.PostWebhook(body); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
