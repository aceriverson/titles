package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"titles.run/services/interfaces"
	"titles.run/titles/models"
)

type AIServiceImpl struct {
	OpenAIKey string
}

type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewAIService() interfaces.AIService {
	return &AIServiceImpl{
		OpenAIKey: os.Getenv("OPENAI_KEY"),
	}
}

func (a *AIServiceImpl) Title(sport string, polygons []models.Polygon, routeMap string, poi []string) (string, error) {
	polygonNames := make([]string, len(polygons))
	for i, polygon := range polygons {
		polygonNames[i] = polygon.Name
	}

	requestBody := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "You are a route-naming assistant. Your task is to create an accurate and concise title for a Strava activity based on the following inputs: \n\n1. A base64-encoded map image depicting the route traveled, drawn in orange, with a green circle indicating the start point and a red circle indicating the endpoint.\n2. The type of activity, which is one of the following: run, ride, etc.\n3. A list of significant points of interest (POIs) the user passed along the route. These POIs may include trails, landmarks, parks, or notable locations.\n4. A list of user-supplied titles to take into consideration when generating the title.\n\nGuidelines for naming the route:\n- Use the points of interest to form a meaningful title. For example, if the activity passes through or near multiple named trails or landmarks, incorporate their names into the title.\n- Prioritize key features that stand out based on the points of interest and the general path shown in the map image. Avoid including generic or irrelevant details.\n- Make the title concise yet descriptive. Aim for 5-8 words.\n- If there is a single prominent landmark, feature, or trail, you can highlight it in the title (e.g., \"Minuteman Trail Run\").\n- If the route involves a loop or is linear, include words like \"Loop\" or \"Point-to-Point\" if applicable.\n- Most importantly, rely on User-Supplied titles. Only return your title with no quotation marks or other words.",
			},
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": "data:image/png;base64," + routeMap,
						},
					},
					{
						"type": "text",
						"text": fmt.Sprintf("Generate a route title for the following activity:\n\n- **Activity Type:** %s\n- **Points of Interest:** %s\n- **Map Description:** The route is drawn in orange on the map, starting at a green circle and ending at a red circle.\n- **User-Supplied Titles:** %s\n\nProvide a concise and descriptive route title based on the inputs.", sport, strings.Join(poi, ", "), strings.Join(polygonNames, ", ")),
					},
				},
			},
		},
		"max_tokens":  300,
		"temperature": 1.3,
	}

	fmt.Println(requestBody)

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
		return "", errors.New("failed to marshal request body")
	}

	url := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return "", errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.OpenAIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return "", errors.New("failed to make request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return "", errors.New("failed to read response body")
	}

	log.Printf("Raw response body: %s", respBody)

	var response CompletionResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Fatalf("Failed to parse response JSON: %v", err)
	}

	fmt.Println(response)

	return response.Choices[0].Message.Content, nil
}
