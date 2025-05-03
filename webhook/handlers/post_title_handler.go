package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	auth "titles.run/jwt"
	"titles.run/webhook/models"
)

func (h *Handler) PostTitleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			log.Println("JWT cookie not found:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			log.Println("Invalid JWT:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var body models.TitleRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := h.titles.PostTitle(userID, body); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
