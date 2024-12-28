package titles

import (
	"titles.run/titles/models"
)

func (h *TitlesCore) PostPolygon(userID string, polygon models.Polygon) error {
	return h.DB.PostPolygon(userID, polygon)
}
