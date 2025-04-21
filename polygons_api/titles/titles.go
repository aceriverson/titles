package titles

import (
	"titles.run/services/interfaces"
)

type TitlesCore struct {
	DB     interfaces.DBService
	Strava interfaces.StravaService
}

func NewTitlesCore(
	db interfaces.DBService,
	strava interfaces.StravaService,
) *TitlesCore {
	core := &TitlesCore{
		DB:     db,
		Strava: strava,
	}

	return core
}
