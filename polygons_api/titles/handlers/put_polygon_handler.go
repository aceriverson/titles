package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"titles.run/auth"
	"titles.run/titles/models"
)

func (h *Handler) PutPolygonHandler() http.HandlerFunc {
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

		var polygon models.Polygon
		err = json.NewDecoder(r.Body).Decode(&polygon)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err = h.titles.PutPolygon(userID, polygon); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
