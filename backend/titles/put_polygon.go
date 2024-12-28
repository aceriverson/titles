package titles

import (
	"fmt"

	"titles.run/titles/models"
)

func (h *TitlesCore) PutPolygon(userID int64, polygon models.Polygon) error {
	fmt.Println(polygon)
	return h.DB.PutPolygon(userID, polygon)
}
