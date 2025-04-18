package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	titlesErrors "titles.run/services/errors"
)

func (h *Handler) GetExchangeTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		scope := r.URL.Query().Get("scope")

		jwt, err := h.titles.GetExchangeToken(code, scope)
		if errors.Is(err, titlesErrors.ErrInvalidScope) {
			http.Error(w, "Invalid scope", http.StatusBadRequest)
			return
		} else if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    jwt,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   259200,
		})

		redirectURL := fmt.Sprintf("https://%s/", os.Getenv("HOST"))
		w.Header().Set("Location", redirectURL)
		w.WriteHeader(http.StatusFound)
	}
}
