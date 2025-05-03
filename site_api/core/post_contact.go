package core

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"titles.run/site_api/models"
)

func (h *Core) PostContact(contact models.Contact) error {
	url := fmt.Sprintf("https://ntfy.sh/%s", os.Getenv("NTFY_CONTACT_CHANNEL"))

	body := fmt.Sprintf("%s - %s - %s", contact.Email, contact.Subject, contact.Message)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Printf("Failed to create ntfy request: %v", err)
		return errors.New("failed to create ntfy request")
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make ntfy request: %v", err)
		return errors.New("failed to make ntfy request")
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to notify of activity: %v", resp.StatusCode)
		return errors.New("failed to notify of activity")
	}

	defer resp.Body.Close()

	return nil
}
