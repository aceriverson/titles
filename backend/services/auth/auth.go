package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateJWT(userID int64) (string, error) {
	claims := &CustomClaims{
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractJWT(r *http.Request) (string, error) {
	jwt := r.Header.Get("Authorization")
	if jwt == "" {
		return "", errors.New("no authorization header found in request")
	}
	return strings.Split(jwt, " ")[1], nil
}

func ValidateJWT(token string) (int64, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return 0, errors.New("JWT_SECRET_KEY not found in environment")
	}

	parsed, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("invalid token or claims")
	}

	if claims, ok := parsed.Claims.(*CustomClaims); ok && parsed.Valid {
		return claims.UserID, nil
	}
	return 0, errors.New("invalid token or claims")
}
