package interfaces

import (
	"titles.run/webhook/models"
)

type AIService interface {
	Title(plan models.UserPlan, activity models.Activity, polygons []models.Polygon, routeMap string, poi []string) (string, error)
}

type DBService interface {
	Close()
	GetIntersectingPolygons(userID int64, points [][]float64) ([]models.Polygon, error)
	GetUserInternal(userID int64) (models.User, error)
	NewUser(user models.User) error
	UnauthorizeUser(userID int64) error
	UpdateUser(user models.User) error
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
	Notify(user models.User, activity models.Activity, update models.Update) error
}

type StravaService interface {
	GetActivity(user models.User, activityID int64) (models.Activity, error)
	RefreshUser(user models.User) (models.User, error)
	RenameActivity(user models.User, activity models.Activity, update models.Update) error
}
