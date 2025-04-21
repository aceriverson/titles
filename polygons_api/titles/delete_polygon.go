package titles

import (
	"titles.run/titles/models"
)

func (h *TitlesCore) DeletePolygon(userID int64, polygon models.Polygon) error {
	return h.DB.DeletePolygon(userID, polygon)
}
