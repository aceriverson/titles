package titles

import (
	"titles.run/titles/models"
)

func (h *TitlesCore) PutPolygon(userID int64, polygon models.Polygon) error {
	return h.DB.PutPolygon(userID, polygon)
}
