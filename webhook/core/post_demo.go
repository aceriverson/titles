package titles

import (
	"errors"

	strava "titles.run/strava/models"
	"titles.run/webhook/models"
)

func (h *TitlesCore) PostDemo(gpx models.Gpx) (string, error) {
	pointsAsSlice := make([][]float64, len(gpx.Points))
	for i, point := range gpx.Points {
		pointsAsSlice[i] = []float64{point.Latitude, point.Longitude}
	}

	routeMap, err := h.Map.GenerateMap(pointsAsSlice)
	if err != nil {
		return "", errors.New("failed to generate map")
	}

	var poi models.POIs
	poi, err = h.DB.GetPOI(pointsAsSlice)

	if err != nil {
		return "", errors.New("failed to get poi via cache")
	}

	titles := make([]string, len(poi.Items))
	for i, item := range poi.Items {
		titles[i] = item.Title
	}

	title, err := h.AI.Title(strava.UserPlanFree, strava.Activity{SportType: "Run"}, nil, routeMap, titles)
	if err != nil {
		return "", errors.New("failed to get title")
	}

	if err := h.Ntfy.Notify(strava.UserInternal{Name: "Demo", ID: 0}, strava.Activity{ID: 0}, strava.Update{Name: title}); err != nil {
		return "", errors.New("failed to send notification")
	}

	return title, nil
}
