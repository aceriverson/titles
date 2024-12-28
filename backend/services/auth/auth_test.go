package auth

import (
	"net/http"
	"os"
	"testing"
)

func TestAuthServiceImpl_CreateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "foo")

	jwt, err := CreateJWT("1234")
	if err != nil {
		t.Fatalf("CreateJWT returned an error: %v", err)
	}

	expectedJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNCJ9._R5CCR7PnbmHC3Kmga1y82qptcn12JAG2ZrhkG1lxqs"
	if jwt != expectedJwt {
		t.Errorf("Expected %s, got %s", expectedJwt, jwt)
	}
}

func TestAuthServiceImpl_ExtractJWT(t *testing.T) {
	req := &http.Request{
		Header: http.Header{
			"Authorization": []string{
				"Bearer foo",
			},
		},
	}

	jwt, err := ExtractJWT(req)

	if err != nil {
		t.Fatalf("ExtractJWT returned an error: %v", err)
	}

	expectedJwt := "foo"
	if jwt != expectedJwt {
		t.Errorf("Expected %s, got %s", expectedJwt, jwt)
	}
}

func TestAuthServiceImpl_ValidateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "foo")

	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNDMyMSJ9.Wu7MM2i52ioZyN8EPNNyYv4psmqWV3JBTyxaZKveOVc"

	userID, err := ValidateJWT(mockToken)

	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	expectedUserID := "4321"
	if userID != expectedUserID {
		t.Errorf("Expected %s, got %s", expectedUserID, userID)
	}
}
