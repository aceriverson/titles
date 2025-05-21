package core

import (
	"log"

	"titles.run/strava/models"
)

func (h *Core) UpdateSettings(userID int64, settings models.Settings) error {
	user, err := h.DB.GetUser(userID)
	if err != nil {
		log.Println("Error getting user:", err)
		return err
	}

	defaultSettings := models.Settings{
		AutomaticTitle: true,
		Tone:           50,
		Attribution:    true,
		Description:    false,
	}

	if user.Plan == "free" {
		defaultSettings.AutomaticTitle = settings.AutomaticTitle
		settings = defaultSettings
	}

	if err := h.DB.UpdateSettings(userID, settings); err != nil {
		return err
	}

	return nil
}
