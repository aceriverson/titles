package handlers

import (
	"net/http"

	titles "titles.run/webhook/core"
)

type Handler struct {
	titles *titles.TitlesCore
}

func RegisterHandlers(titles *titles.TitlesCore) {
	handler := &Handler{titles}

	http.Handle("GET /webhook", handler.GetWebhookHandler())
	http.Handle("POST /webhook", handler.PostWebhookHandler())
}
