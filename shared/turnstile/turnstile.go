package turnstile

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type TurnStile struct {
	Secret string
	Token  string
}

type TurnStileResponse struct {
	Success bool `json:"success"`
}

func (t *TurnStile) Verify() error {
	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	requestBody := map[string]string{
		"secret":   t.Secret,
		"response": t.Token,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tsResponse TurnStileResponse
	if err := json.NewDecoder(resp.Body).Decode(&tsResponse); err != nil {
		return err
	}

	if !tsResponse.Success {
		return errors.New("turnstile verification failed")
	}

	return nil
}
