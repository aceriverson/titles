package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"titles.run/services/auth"
	"titles.run/titles/models"
)

func (h *Handler) DeletePolygonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwt, err := auth.ExtractJWT(r)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(jwt)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var polygon models.Polygon
		err = json.NewDecoder(r.Body).Decode(&polygon)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := h.titles.DeletePolygon(userID, polygon); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
