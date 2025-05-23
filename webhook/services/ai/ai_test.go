package ai

import (
	"os"
	"testing"

	strava "titles.run/strava/models"

	"github.com/jarcoal/httpmock"
	"titles.run/webhook/models"
)

func TestAIServiceImpl_Title(t *testing.T) {
	os.Setenv("AI_URL", "https://ai.example.com/v1/chat/completions")
	os.Setenv("AI_API_KEY", "dummy_key")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockResponse := `{
		"choices": [
			{
				"message": {
					"content": "Beautiful Mountain Trail Loop"
				}
			}
		]
	}`
	httpmock.RegisterResponder(
		"POST",
		"https://ai.example.com/v1/chat/completions",
		httpmock.NewStringResponder(200, mockResponse),
	)

	service := NewAIService()

	activity := strava.Activity{
		SportType: "run",
		SegmentEfforts: []strava.SegmentEffort{
			{Name: "Central Park Loop"},
			{Name: "Central Park East"},
		},
	}
	polygons := []models.Polygon{
		{Name: "Central Park"},
		{Name: "Harlem River Greenway"},
	}
	routeMap := "data:image/png;base64,FAKE_BASE64_DATA"
	poi := []string{"Central Park", "Harlem River"}

	title, err := service.Title(strava.UserPlanFree, 50, activity, polygons, routeMap, poi)
	if err != nil {
		t.Fatalf("Title method returned an error: %v", err)
	}

	expectedTitle := "Beautiful Mountain Trail Loop"
	if title != expectedTitle {
		t.Errorf("Expected title '%s', got '%s'", expectedTitle, title)
	}

	info := httpmock.GetCallCountInfo()
	if count := info["POST https://ai.example.com/v1/chat/completions"]; count != 1 {
		t.Errorf("Expected 1 call to the mocked endpoint, but got %d", count)
	}
}
