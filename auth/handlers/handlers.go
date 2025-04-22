package handlers

import (
	"net/http"

	auth "titles.run/auth/core"
)

type Handler struct {
	auth *auth.Core
}

func RegisterHandlers(auth *auth.Core) {
	handler := &Handler{auth}

	http.Handle("GET /callback", handler.GetCallbackHandler())
	http.Handle("POST /logout", handler.PostLogoutHandler())
}
