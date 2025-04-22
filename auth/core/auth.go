package auth

import (
	"titles.run/auth/services/interfaces"
	"titles.run/strava"
)

type Core struct {
	DB     interfaces.DBService
	Strava strava.Client
}

func NewAuthCore(db interfaces.DBService, strava strava.Client) *Core {
	core := &Core{
		DB:     db,
		Strava: strava,
	}

	return core
}
