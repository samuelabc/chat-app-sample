package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token for a user
func GenerateJWT(userId int, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the username
func ValidateJWT(tokenString string) (int, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, "", errors.New("invalid token signature")
		}
		return 0, "", err
	}
	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}
	return claims.UserId, claims.Username, nil
}
