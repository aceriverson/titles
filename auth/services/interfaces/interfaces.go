package interfaces

import "titles.run/strava/models"

type DBService interface {
	Close()
	NewUser(user models.UserInternal) error
}
