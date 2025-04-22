package main

import (
	"log"
	"net/http"

	"os"

	auth "titles.run/auth/core"
	"titles.run/auth/handlers"
	"titles.run/auth/services/db"
	"titles.run/strava"
)

func main() {
	log.Println("Starting auth service...")

	dbService := db.NewDBService()
	stravaService := strava.NewStravaService(os.Getenv("STRAVA_CLIENT_ID"), os.Getenv("STRAVA_CLIENT_SECRET"))

	defer dbService.Close()

	if dbService != nil && stravaService != nil {
		log.Println("All services initialized successfully.")
	} else {
		log.Println("One or more services failed to initialize.")
	}

	core := auth.NewAuthCore(dbService, stravaService)

	handlers.RegisterHandlers(core)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
