package main

import (
	"log"
	"net/http"

	site "titles.run/site_api/core"
	"titles.run/site_api/handlers"
	"titles.run/site_api/services/db"
)

func main() {
	log.Println("Starting site_api service...")

	dbService := db.NewDBService()

	defer dbService.Close()

	if dbService != nil {
		log.Println("All services initialized successfully.")
	} else {
		log.Println("One or more services failed to initialize.")
	}

	core := site.NewSiteCore(dbService)

	handlers.RegisterHandlers(core)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
