package main

import (
	"log"
	"net/http"

	"titles.run/webhook/services/ai"
	"titles.run/webhook/services/db"
	"titles.run/webhook/services/dedupe"
	"titles.run/webhook/services/here"
	map_service "titles.run/webhook/services/map"
	"titles.run/webhook/services/ntfy"
	"titles.run/webhook/services/strava"

	titles "titles.run/webhook/core"
	"titles.run/webhook/handlers"
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
