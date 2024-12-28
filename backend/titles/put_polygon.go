package titles

import (
	"fmt"

	"titles.run/titles/models"
)

func (h *TitlesCore) PutPolygon(userID string, polygon models.Polygon) error {
	fmt.Println(polygon)
	return h.DB.PutPolygon(userID, polygon)
}
