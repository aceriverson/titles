package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

		redirectURL := fmt.Sprintf("https://%s/?token=%s", os.Getenv("HOST"), url.QueryEscape(jwt))
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
