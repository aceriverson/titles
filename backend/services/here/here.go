package here

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"titles.run/services/interfaces"
)

type HereServiceImpl struct {
	HereAppID  string
	HereAPIKey string
}

type APIResponse struct {
	Items []struct {
		Title string `json:"title"`
	} `json:"items"`
}

func NewHereService() interfaces.HereService {
	return &HereServiceImpl{
		HereAppID:  os.Getenv("HERE_APP_ID"),
		HereAPIKey: os.Getenv("HERE_API_KEY"),
	}
}

func (h *HereServiceImpl) GetPOI(line string, at []float64) ([]string, error) {
	url := fmt.Sprintf("https://browse.search.hereapi.com/v1/browse?route=%s;w=50&limit=10&categories=300-3000-0023,300-3000-0025,300-3000-0030,300-3000-0450,350-3500,350-3510,350-3522,350-3550,550-5510,800-8600-0195,800-8600-0381&apiKey=%s&at=%f,%f", line, h.HereAPIKey, at[0], at[1])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return []string{}, errors.New("failed to create HERE request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return []string{}, errors.New("failed to make HERE request")
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return []string{}, errors.New("failed to read HERE response body")
	}

	var response APIResponse
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
		return []string{}, errors.New("failed to unmarshal HERE response")
	}

	titles := make([]string, len(response.Items))
	for i, item := range response.Items {
		titles[i] = item.Title
	}

	return titles, nil
}
