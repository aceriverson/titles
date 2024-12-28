package interfaces

import (
	"titles.run/titles/models"
)

type AIService interface {
	Title(sport string, polygons []models.Polygon, routeMap string, poi []string) (string, error)
}

type DBService interface {
	Close()
	DeletePolygon(userID string, polygon models.Polygon) error
	GetIntersectingPolygons(userID string, points [][]float64) ([]models.Polygon, error)
	GetUser(userID string) (models.User, error)
	GetUserInternal(userID string) (models.UserInternal, error)
	GetPolygons(userID string) ([]models.Polygon, error)
	NewUser(user models.UserInternal) error
	PostPolygon(userID string, polygon models.Polygon) error
	PutPolygon(userID string, polygon models.Polygon) error
	UnauthorizeUser(userID string) error
	UpdateUser(user models.UserInternal) error
}

type DedupeService interface {
	Close()
	AddActivity(id string) error
	DedupeActivity(id string) (bool, error)
}

type HereService interface {
	GetPOI(line string, start []float64) ([]string, error)
}

type MapService interface {
	GenerateMap(coords [][]float64) (string, error)
}

type NtfyService interface {
	Notify(user models.UserInternal, activity models.Activity, update models.Update) error
}

type StravaService interface {
	GetActivity(user models.UserInternal, activityID string) (models.Activity, error)
	RefreshUser(user models.UserInternal) (models.UserInternal, error)
	RenameActivity(user models.UserInternal, activity models.Activity, update models.Update) error
	TokenExchange(code string) (models.UserInternal, error)
}
