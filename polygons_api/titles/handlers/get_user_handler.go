package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	auth "titles.run/jwt"
)

func (h *Handler) GetUserHandler() http.HandlerFunc {
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
