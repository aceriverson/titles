package handlers

import (
	"net/http"

	site "titles.run/site_api/core"
)

type Handler struct {
	core *site.Core
}

func RegisterHandlers(core *site.Core) {
	handler := &Handler{core}

	http.Handle("GET /user", handler.GetUserHandler())

	http.Handle("POST /accept_terms", handler.PostAcceptTermsHandler())
	http.Handle("POST /contact", handler.PostContactHandler())
}
