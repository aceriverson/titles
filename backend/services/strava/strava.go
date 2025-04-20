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

	"titles.run/services/interfaces"
	"titles.run/titles/models"
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

func (s *StravaServiceImpl) RefreshUser(user models.UserInternal) (models.UserInternal, error) {
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
		return models.UserInternal{}, errors.New("failed to marshal request body")
	}

	url := "https://www.strava.com/api/v3/oauth/token"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return models.UserInternal{}, errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return models.UserInternal{}, errors.New("failed to make request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return models.UserInternal{}, errors.New("failed to read response body")
	}

	var refreshedToken models.RefreshResponse
	if err := json.Unmarshal(respBody, &refreshedToken); err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
		return models.UserInternal{}, err
	}

	user.UpdateToken(refreshedToken)

	return user, nil
}

func (s *StravaServiceImpl) TokenExchange(code string) (models.UserInternal, error) {
	tokenExchangeURL := fmt.Sprintf(
		"https://www.strava.com/api/v3/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code",
		s.ClientID,
		s.ClientSecret,
		code,
	)
	resp, err := http.Post(tokenExchangeURL, "application/json", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.UserInternal{}, errors.New("unable to exchange code for token")
	}
	defer resp.Body.Close()

	var oauthResponse models.TokenExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResponse); err != nil {
		return models.UserInternal{}, err
	}

	return oauthResponse.ToUserInternal(), nil
}
