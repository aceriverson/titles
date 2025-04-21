package ntfy

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"titles.run/webhook/models"
	"titles.run/webhook/services/interfaces"
)

type NtfyServiceImpl struct {
	NtfyChannel string
}

func NewNtfyService() interfaces.NtfyService {
	return &NtfyServiceImpl{
		NtfyChannel: os.Getenv("NTFY_CHANNEL"),
	}
}

func (n *NtfyServiceImpl) Notify(user models.User, activity models.Activity, update models.Update) error {
	url := fmt.Sprintf("https://ntfy.sh/%s", n.NtfyChannel)

	body := fmt.Sprintf("%s - %d - %s", user.Name, activity.ID, update.Name)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return errors.New("failed to make request")
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to notify of activity: %v", resp.StatusCode)
		return errors.New("failed to notify of activity")
	}

	defer resp.Body.Close()

	return nil
}
