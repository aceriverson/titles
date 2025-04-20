package titles

import (
	"titles.run/webhook/services/interfaces"
)

type TitlesCore struct {
	DB     interfaces.DBService
	Dedupe interfaces.DedupeService
	AI     interfaces.AIService
	Strava interfaces.StravaService
	Map    interfaces.MapService
	Here   interfaces.HereService
	Ntfy   interfaces.NtfyService
}

func NewTitlesCore(
	db interfaces.DBService,
	dedupe interfaces.DedupeService,
	ai interfaces.AIService,
	strava interfaces.StravaService,
	map_service interfaces.MapService,
	here interfaces.HereService,
	ntfy interfaces.NtfyService,
) *TitlesCore {
	core := &TitlesCore{
		DB:     db,
		Dedupe: dedupe,
		AI:     ai,
		Strava: strava,
		Map:    map_service,
		Here:   here,
		Ntfy:   ntfy,
	}

	return core
}
