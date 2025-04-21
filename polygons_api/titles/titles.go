package titles

import (
	"titles.run/services/interfaces"
	"titles.run/strava"
)

type TitlesCore struct {
	DB     interfaces.DBService
	Strava strava.Client
}

func NewTitlesCore(
	db interfaces.DBService,
	strava strava.Client,
) *TitlesCore {
	core := &TitlesCore{
		DB:     db,
		Strava: strava,
	}

	return core
}
