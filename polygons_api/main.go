package main

import (
	"log"
	"net/http"
	"os"

	"titles.run/services/db"
	"titles.run/strava"

	"titles.run/titles"
	"titles.run/titles/handlers"
)

func main() {
	log.Println("Starting application...")

	dbService := db.NewDBService()
	stravaService := strava.NewStravaService(os.Getenv("STRAVA_CLIENT_ID"), os.Getenv("STRAVA_CLIENT_SECRET"))

	defer dbService.Close()

	if dbService != nil && stravaService != nil {
		log.Println("All services initialized successfully.")
	} else {
		log.Println("One or more services failed to initialize.")
	}

	titlesCore := titles.NewTitlesCore(dbService, stravaService)

	handlers.RegisterHandlers(titlesCore)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
