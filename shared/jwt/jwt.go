package jwt

import (
	"errors"
	"os"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwtLib.RegisteredClaims
}

func CreateJWT(userID int64) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(expirationTime),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(token string) (int64, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return 0, errors.New("JWT_SECRET_KEY not found in environment")
	}

	parsed, err := jwtLib.ParseWithClaims(token, &CustomClaims{}, func(token *jwtLib.Token) (interface{}, error) {
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
