package titles

import (
	"titles.run/titles/models"
)

func (h *TitlesCore) GetUser(userID string) (models.User, error) {
	return h.DB.GetUser(userID)
}
