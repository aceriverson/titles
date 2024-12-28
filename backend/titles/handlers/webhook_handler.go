package handlers

import "net/http"

func (h *Handler) WebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
