package titles

import (
	"titles.run/titles/models"
)

func (h *TitlesCore) GetPolygons(userID int64) ([]models.Polygon, error) {
	return h.DB.GetPolygons(userID)
}
