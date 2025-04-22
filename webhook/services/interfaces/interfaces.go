package interfaces

import (
	strava "titles.run/strava/models"
	"titles.run/webhook/models"
)

type AIService interface {
	Title(plan strava.UserPlan, activity strava.Activity, polygons []models.Polygon, routeMap string, poi []string) (string, error)
}

type DBService interface {
	Close()
	GetIntersectingPolygons(userID int64, points [][]float64) ([]models.Polygon, error)
	GetPOI(points [][]float64) (models.POIs, error)
	GetUserInternal(userID int64) (strava.UserInternal, error)
	SetPOI(poi models.POIs) error
	UnauthorizeUser(userID int64) error
	UpdateUser(user strava.UserInternal) error
}

type TTLStoreService interface {
	Close()
	AddActivity(id int64) error
	DedupeActivity(id int64) (bool, error)
	CheckRateLimit(userID int64, plan strava.UserPlan) (bool, error)
	IncrementRateLimit(userID int64) error
}

type HereService interface {
	GetPOI(line string, start []float64) (models.POIs, error)
}

type MapService interface {
	GenerateMap(coords [][]float64) (string, error)
}

type NtfyService interface {
	Notify(user strava.UserInternal, activity strava.Activity, update strava.Update) error
}
