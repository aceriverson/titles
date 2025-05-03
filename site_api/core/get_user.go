package core

import (
	"titles.run/strava/models"
)

func (h *Core) GetUser(userID int64) (models.User, error) {
	return h.DB.GetUser(userID)
}
