package interfaces

import "titles.run/strava/models"

type DBService interface {
	Close()
	AcceptTerms(userID int64) error
	GetUser(userID int64) (models.User, error)
	GetUserInternal(userID int64) (models.UserInternal, error)
}
