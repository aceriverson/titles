package auth

import (
	"titles.run/auth/errors"
	"titles.run/jwt"
)

func (h *Core) GetExchangeToken(code, scope string) (string, error) {
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

	return jwt.CreateJWT(user.ID)
}
