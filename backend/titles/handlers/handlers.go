package handlers

import (
	"net/http"

	"titles.run/titles"
)

type Handler struct {
	titles *titles.TitlesCore
}

func RegisterHandlers(titles *titles.TitlesCore) {
	handler := &Handler{titles}

	http.Handle("GET /exchange_token", handler.GetExchangeTokenHandler())

	http.Handle("/webhook", handler.WebhookHandler())

	http.Handle("DELETE /polygon", handler.DeletePolygonHandler())
	http.Handle("POST /polygon", handler.PostPolygonHandler())
	http.Handle("PUT /polygon", handler.PutPolygonHandler())

	http.Handle("GET /polygons", handler.GetPolygonsHandler())

	http.Handle("GET /user", handler.GetUserHandler())
	http.Handle("POST /logout", handler.PostLogoutHandler())

	http.Handle("GET /webhook", handler.GetWebhookHandler())
	http.Handle("POST /webhook", handler.PostWebhookHandler())
}
