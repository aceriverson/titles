package interfaces

import (
	"titles.run/titles/models"
)

type AIService interface {
	Title(sport string, polygons []models.Polygon, routeMap string, poi []string) (string, error)
}

type DBService interface {
	Close()
	DeletePolygon(userID int64, polygon models.Polygon) error
	GetIntersectingPolygons(userID int64, points [][]float64) ([]models.Polygon, error)
	GetUser(userID int64) (models.User, error)
	GetUserInternal(userID int64) (models.UserInternal, error)
	GetPolygons(userID int64) ([]models.Polygon, error)
	NewUser(user models.UserInternal) error
	PostPolygon(userID int64, polygon models.Polygon) error
	PutPolygon(userID int64, polygon models.Polygon) error
	UnauthorizeUser(userID int64) error
	UpdateUser(user models.UserInternal) error
}

type DedupeService interface {
	Close()
	AddActivity(id int64) error
	DedupeActivity(id int64) (bool, error)
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
	GetActivity(user models.UserInternal, activityID int64) (models.Activity, error)
	RefreshUser(user models.UserInternal) (models.UserInternal, error)
	RenameActivity(user models.UserInternal, activity models.Activity, update models.Update) error
	TokenExchange(code string) (models.UserInternal, error)
}
