package titles

import (
	"titles.run/titles/models"
)

func (h *TitlesCore) GetUser(userID int64) (models.User, error) {
	return h.DB.GetUser(userID)
}
