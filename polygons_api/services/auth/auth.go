package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateJWT(userID int64) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	return token.SignedString([]byte(jwtSecret))
}

func ExtractJWT(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return "", errors.New("jwt cookie not found")
	}
	return cookie.Value, nil
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
