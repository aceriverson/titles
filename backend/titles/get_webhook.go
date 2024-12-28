package titles

import (
	"errors"
	"fmt"
	"os"
)

func (h *TitlesCore) GetWebhook(verify_token string) error {
	fmt.Println("verify_token:", verify_token)
	if verify_token == os.Getenv("STRAVA_WEBHOOK_VERIFY_TOKEN") {
		return nil
	}

	return errors.New("invalid verification token")
}
