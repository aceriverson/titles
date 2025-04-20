package titles

import (
	"errors"
	"os"
)

func (h *TitlesCore) GetWebhook(verify_token string) error {
	if verify_token == os.Getenv("STRAVA_WEBHOOK_VERIFY_TOKEN") {
		return nil
	}

	return errors.New("invalid verification token")
}
