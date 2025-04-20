package interfaces

import (
	"titles.run/titles/models"
)

type DBService interface {
	Close()
	DeletePolygon(userID int64, polygon models.Polygon) error
	GetUser(userID int64) (models.User, error)
	GetUserInternal(userID int64) (models.UserInternal, error)
	GetPolygons(userID int64) ([]models.Polygon, error)
	NewUser(user models.UserInternal) error
	PostPolygon(userID int64, polygon models.Polygon) error
	PutPolygon(userID int64, polygon models.Polygon) error
}

type StravaService interface {
	RefreshUser(user models.UserInternal) (models.UserInternal, error)
	TokenExchange(code string) (models.UserInternal, error)
}
