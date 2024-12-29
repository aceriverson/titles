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
				"role": "system",
				"content": `
					You are a route-naming assistant for Strava activities. Your primary goal is to generate an accurate and concise title for an activity based on the following inputs:\n\n
					
					1. A base64-encoded map image depicting the route traveled, drawn in orange, with a green circle indicating the starting point and a red circle indicating the endpoint.\n
					2. The type of activity, such as run, ride, etc.\n
					3. A list of significant points of interest (POIs) the user passed along the route. These may include trails, landmarks, parks, or notable locations.\n
					4. A list of user-supplied titles to take into consideration when generating the title.\n\n
					
					### Guidelines for Naming the Route:\n
					1. **Focus on Key Features**:\n
					   - Prioritize naming the route based on significant features encountered farthest from the starting point, such as notable trails, lakes, or parks.\n
					   - If the route prominently follows a trail, highlight the trail name (e.g., \"Charles River Trail\").\n
					   - If the route circles or traverses a prominent feature, such as a lake or park, include it in the title (e.g., \"Lake Loop Run\").\n\n
					   
					2. **Incorporate Points of Interest**:\n
					   - Use POIs to form a meaningful and descriptive title.\n
					   - Prioritize features based on prominence or uniqueness, with a focus on features farthest from the start.\n
					   - Avoid generic or irrelevant details (e.g., \"Neighborhood Ride\").\n\n
					   
					3. **Keep it Concise**:\n
					   - Aim for a title between 5-8 words that captures the essence of the activity.\n
					   - Include terms like \"Loop\" or \"Out-and-Back\" if applicable.\n\n
					   
					4. **Use User-Supplied Titles**:\n
					   - If the user provides a title, ensure it influences the generated title.\n
					   - Blend user-supplied titles with your insights if necessary.\n\n
					   
					### Important Note:\n
					- Rely on the provided map image and POIs to identify the most defining aspects of the route.\n
					- Return only the generated title without any additional comments, quotation marks, or formatting.`,
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
						"text": fmt.Sprintf(`
							Generate a route title for the following activity:\n\n
							
							- **Activity Type:** %s\n
							- **Points of Interest:** %s\n
							- **Map Description:** The route is drawn in orange on the map, starting at a green circle and ending at a red circle.\n
							- **User-Supplied Titles:** %s\n\n
							
							Provide a concise and descriptive route title based on the inputs.`,
							sport, strings.Join(poi, ", "), strings.Join(polygonNames, ", "),
						),
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
