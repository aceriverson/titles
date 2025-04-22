package interfaces

import (
	strava "titles.run/strava/models"
	"titles.run/titles/models"
)

type DBService interface {
	Close()
	DeletePolygon(userID int64, polygon models.Polygon) error
	GetUser(userID int64) (strava.User, error)
	GetPolygons(userID int64) ([]models.Polygon, error)
	PostPolygon(userID int64, polygon models.Polygon) error
	PutPolygon(userID int64, polygon models.Polygon) error
}
