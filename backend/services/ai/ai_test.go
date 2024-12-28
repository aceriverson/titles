package ai

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"titles.run/titles/models"
)

func TestAIServiceImpl_Title(t *testing.T) {
	os.Setenv("OPENAI_KEY", "dummy_key")

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
		"https://api.openai.com/v1/chat/completions",
		httpmock.NewStringResponder(200, mockResponse),
	)

	service := NewAIService()

	sport := "run"
	polygons := []models.Polygon{
		{Name: "Central Park"},
		{Name: "Harlem River Greenway"},
	}
	routeMap := "data:image/png;base64,FAKE_BASE64_DATA"
	poi := []string{"Central Park", "Harlem River"}

	title, err := service.Title(sport, polygons, routeMap, poi)
	if err != nil {
		t.Fatalf("Title method returned an error: %v", err)
	}

	expectedTitle := "Beautiful Mountain Trail Loop"
	if title != expectedTitle {
		t.Errorf("Expected title '%s', got '%s'", expectedTitle, title)
	}

	info := httpmock.GetCallCountInfo()
	if count := info["POST https://api.openai.com/v1/chat/completions"]; count != 1 {
		t.Errorf("Expected 1 call to the mocked endpoint, but got %d", count)
	}
}
