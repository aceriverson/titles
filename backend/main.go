package main

import (
	"log"
	"net/http"

	"titles.run/services/ai"
	"titles.run/services/db"
	"titles.run/services/dedupe"
	"titles.run/services/here"
	map_service "titles.run/services/map"
	"titles.run/services/ntfy"
	"titles.run/services/strava"

	"titles.run/titles"
	"titles.run/titles/handlers"
)

func main() {
	log.Println("Starting application...")

	dbService := db.NewDBService()
	dedupeService := dedupe.NewDedupeService()
	aiService := ai.NewAIService()
	stravaService := strava.NewStravaService()
	mapService := map_service.NewMapService()
	hereService := here.NewHereService()
	ntfyService := ntfy.NewNtfyService()

	defer dbService.Close()
	defer dedupeService.Close()

	if dbService != nil && dedupeService != nil && aiService != nil && stravaService != nil && mapService != nil && hereService != nil && ntfyService != nil {
		log.Println("All services initialized successfully.")
	} else {
		log.Println("One or more services failed to initialize.")
	}

	titlesCore := titles.NewTitlesCore(dbService, dedupeService, aiService, stravaService, mapService, hereService, ntfyService)

	handlers.RegisterHandlers(titlesCore)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
