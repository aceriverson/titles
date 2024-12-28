package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"titles.run/services/auth"
)

func (h *Handler) GetUserHandler() http.HandlerFunc {
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

		userData, err := h.titles.GetUser(userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			log.Println(err)
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	}
}
