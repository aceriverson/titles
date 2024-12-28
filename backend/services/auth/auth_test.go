package auth

import (
	"net/http"
	"os"
	"testing"
)

func TestAuthServiceImpl_CreateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "foo")

	jwt, err := CreateJWT(1234)
	if err != nil {
		t.Fatalf("CreateJWT returned an error: %v", err)
	}

	expectedJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjM0fQ.dAFL0yHMYCqdNvLjmRkNG-f3rX4HqVCEsUa36g_jt1s"
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

	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0MzIxfQ.glFdbMEN4gITHdrRO9MUTn2dIhp04CimMU97w1EQ8Dw"

	userID, err := ValidateJWT(mockToken)

	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	expectedUserID := int64(4321)
	if userID != expectedUserID {
		t.Errorf("Expected %d, got %d", expectedUserID, userID)
	}
}
