package main

import (
	"log"
	"net/http"
	"os"

	"titles.run/strava"
	"titles.run/webhook/services/ai"
	"titles.run/webhook/services/db"
	"titles.run/webhook/services/here"
	map_service "titles.run/webhook/services/map"
	"titles.run/webhook/services/ntfy"
	"titles.run/webhook/services/ttlstore"

	titles "titles.run/webhook/core"
	"titles.run/webhook/handlers"
)

func main() {
	log.Println("Starting application...")

	dbService := db.NewDBService()
	ttlStoreService := ttlstore.NewTTLStoreService()
	aiService := ai.NewAIService()
	stravaService := strava.NewStravaService(os.Getenv("STRAVA_CLIENT_ID"), os.Getenv("STRAVA_CLIENT_SECRET"))
	mapService := map_service.NewMapService()
	hereService := here.NewHereService()
	ntfyService := ntfy.NewNtfyService()

	defer dbService.Close()
	defer ttlStoreService.Close()

	if dbService != nil && ttlStoreService != nil && aiService != nil && stravaService != nil && mapService != nil && hereService != nil && ntfyService != nil {
		log.Println("All services initialized successfully.")
	} else {
		log.Println("One or more services failed to initialize.")
	}

	titlesCore := titles.NewTitlesCore(dbService, ttlStoreService, aiService, stravaService, mapService, hereService, ntfyService)

	handlers.RegisterHandlers(titlesCore)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
