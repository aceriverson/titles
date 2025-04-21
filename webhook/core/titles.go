package titles

import (
	"titles.run/strava"
	"titles.run/webhook/services/interfaces"
)

type TitlesCore struct {
	DB       interfaces.DBService
	TTLStore interfaces.TTLStoreService
	AI       interfaces.AIService
	Strava   strava.Client
	Map      interfaces.MapService
	Here     interfaces.HereService
	Ntfy     interfaces.NtfyService
}

func NewTitlesCore(
	db interfaces.DBService,
	ttlStore interfaces.TTLStoreService,
	ai interfaces.AIService,
	strava strava.Client,
	mapService interfaces.MapService,
	here interfaces.HereService,
	ntfy interfaces.NtfyService,
) *TitlesCore {
	core := &TitlesCore{
		DB:       db,
		TTLStore: ttlStore,
		AI:       ai,
		Strava:   strava,
		Map:      mapService,
		Here:     here,
		Ntfy:     ntfy,
	}

	return core
}
