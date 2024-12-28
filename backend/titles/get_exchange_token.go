package titles

import (
	"titles.run/services/auth"
	"titles.run/services/errors"
)

func (h *TitlesCore) GetExchangeToken(code, scope string) (string, error) {
	if scope != "read,activity:write,activity:read_all" {
		return "", errors.ErrInvalidScope
	}

	user, err := h.Strava.TokenExchange(code)
	if err != nil {
		return "", err
	}

	if err := h.DB.NewUser(user); err != nil {
		return "", err
	}

	return auth.CreateJWT(user.ID)
}
