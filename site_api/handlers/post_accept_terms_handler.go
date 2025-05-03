package handlers

import (
	"log"
	"net/http"

	auth "titles.run/jwt"
)

func (h *Handler) PostAcceptTermsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			log.Println("Invalid JWT:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := h.core.AcceptTerms(userID); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
