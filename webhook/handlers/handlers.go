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

	http.Handle("GET /", handler.GetWebhookHandler())
	http.Handle("POST /", handler.PostWebhookHandler())

	http.Handle("POST /demo", handler.PostDemoHandler())
	http.Handle("POST /title", handler.PostTitleHandler())

	http.Handle("POST /stripe", handler.PostStripeHandler())
	http.Handle("POST /stripe/", handler.PostStripeHandler())
}
