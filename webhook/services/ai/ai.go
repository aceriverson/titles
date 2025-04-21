package ai

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"titles.run/webhook/models"
	"titles.run/webhook/services/interfaces"
)

//go:embed templates/system_prompt.tmpl
var systemPrompt string

//go:embed templates/user_prompt.tmpl
var userPrompt string

type AIServiceImpl struct {
	AIURL    string
	AIAPIKey string
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
		AIURL:    os.Getenv("AI_URL"),
		AIAPIKey: os.Getenv("AI_API_KEY"),
	}
}

func (a *AIServiceImpl) Title(plan models.UserPlan, activity models.Activity, polygons []models.Polygon, routeMap string, poi []string) (string, error) {
	requestBody, err := constructRequestBody(plan, polygons, activity, routeMap, poi)
	if err != nil {
		log.Printf("Error constructing request body: %v", err)
		return "", err
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
		return "", errors.New("failed to marshal request body")
	}

	url := a.AIURL
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return "", errors.New("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.AIAPIKey))

	resp, err := http.DefaultClient.Do(req)
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

	var response CompletionResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Printf("Failed to parse response JSON: %v", err)
		return "", errors.New("failed to parse response JSON")
	}

	return response.Choices[0].Message.Content, nil
}

func renderTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func constructRequestBody(plan models.UserPlan, polygons []models.Polygon, activity models.Activity, routeMap string, poi []string) (map[string]interface{}, error) {
	polygonNames := make([]string, len(polygons))
	for i, polygon := range polygons {
		polygonNames[i] = polygon.Name
	}

	segmentEfforts := make([]string, len(activity.SegmentEfforts))
	for i, segment := range activity.SegmentEfforts {
		segmentEfforts[i] = segment.Name
	}

	systemContent, err := renderTemplate(systemPrompt, nil)
	if err != nil {
		log.Printf("Error rendering system template: %v", err)
		return map[string]interface{}{}, err
	}

	userContent, err := renderTemplate(userPrompt, map[string]interface{}{
		"ActivityType": activity.SportType,
		"SegmentNames": strings.Join(segmentEfforts, ", "),
		"POIs":         strings.Join(poi, ", "),
		"UserTitles":   strings.Join(polygonNames, ", "),
	})
	if err != nil {
		log.Printf("Error rendering user template: %v", err)
		return map[string]interface{}{}, err
	}

	// The following models are available options for future use if neccessary.
	// "meta-llama/llama-4-scout:free",
	// "qwen/qwen2.5-vl-3b-instruct:free",
	// "qwen/qwen2.5-vl-32b-instruct:free",
	// "mistralai/mistral-small-3.1-24b-instruct:free",
	// "google/gemma-3-1b-it:free",
	// "google/gemma-3-12b-it:free",
	// "qwen/qwen2.5-vl-72b-instruct:free",
	// "qwen/qwen-2.5-vl-7b-instruct:free",
	// "google/gemini-flash-1.5-8b-exp",

	var model string
	if plan == models.UserPlanFree {
		model = "google/gemini-2.0-flash-exp:free"
	} else if plan == models.UserPlanPro {
		model = "google/gemini-2.0-flash-exp:free"
	} else {
		model = "google/gemini-2.0-flash-exp:free"
	}

	var extraModels []string
	if plan == models.UserPlanFree {
		extraModels = []string{
			"google/gemma-3-27b-it:free",
			"meta-llama/llama-4-maverick:free",
			"google/gemma-3-4b-it",
		}
	} else if plan == models.UserPlanPro {
		extraModels = []string{
			"google/gemini-2.0-flash-001",
			"google/gemma-3-27b-it:free",
			"meta-llama/llama-4-maverick:free",
			"google/gemma-3-4b-it",
		}
	} else {
		extraModels = []string{}
	}

	requestBody := map[string]interface{}{
		"model": model,
		"extra_body": map[string]interface{}{
			"models": extraModels,
		},
		"messages": []map[string]interface{}{
			{"role": "system", "content": systemContent},
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
						"text": userContent,
					},
				},
			},
		},
		"max_tokens":  300,
		"temperature": 1.3,
	}

	return requestBody, nil
}
