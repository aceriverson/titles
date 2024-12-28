package titles

import (
	"titles.run/titles/models"
)

func (h *TitlesCore) DeletePolygon(userID string, polygon models.Polygon) error {
	return h.DB.DeletePolygon(userID, polygon)
}
