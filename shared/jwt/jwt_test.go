package jwt

import (
	"os"
	"testing"
)

func TestJWTServiceImpl_CreateJWT(t *testing.T) {
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

func TestJWTServiceImpl_ValidateJWT(t *testing.T) {
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
