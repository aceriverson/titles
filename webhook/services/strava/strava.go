package strava

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"titles.run/webhook/models"
	"titles.run/webhook/services/interfaces"
)

type StravaServiceImpl struct {
	ClientID     string
	ClientSecret string
}

func NewStravaService() interfaces.StravaService {
	return &StravaServiceImpl{
		ClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
	}
}

func (s *StravaServiceImpl) GetActivity(user models.User, activityID int64) (models.Activity, error) {
	url := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d", activityID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return models.Activity{}, errors.New("failed to create request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return models.Activity{}, errors.New("failed to make request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return models.Activity{}, errors.New("failed to read response body")
	}

	var activity models.Activity
	if err := json.Unmarshal(respBody, &activity); err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
		return models.Activity{}, err
	}

	return activity, nil
}

func (s *StravaServiceImpl) RefreshUser(user models.User) (models.User, error) {
	if user.ExpiresAt >= time.Now().Unix() {
		return user, nil
	}

	body := models.RefreshBody{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: user.RefreshToken,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
		return models.User{}, errors.New("failed to marshal request body")
	}

	url := "https://www.strava.com/api/v3/oauth/token"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return models.User{}, errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return models.User{}, errors.New("failed to make request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return models.User{}, errors.New("failed to read response body")
	}

	var refreshedToken models.RefreshResponse
	if err := json.Unmarshal(respBody, &refreshedToken); err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
		return models.User{}, err
	}

	user.UpdateToken(refreshedToken)

	return user, nil
}

func (s *StravaServiceImpl) RenameActivity(user models.User, activity models.Activity, update models.Update) error {
	url := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d", activity.ID)

	body, err := json.Marshal(update)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
		return errors.New("failed to marshal request body")
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return errors.New("failed to make request")
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to rename activity: %v", resp.StatusCode)
		return errors.New("failed to rename activity")
	}

	defer resp.Body.Close()

	return nil
}
