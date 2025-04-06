package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jimvid/dionysus/internal/model"
)

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secret, err := GetJWTSecret()
	if err != nil {
		return jwt.MapClaims{}, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token of unrecognized type - unauthorized")
	}

	return claims, nil
}

func CreateToken(user model.User) (string, error) {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()
	secret, err := GetJWTSecret()

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetJWTSecret() (string, error) {
	secret := "secret"

	return secret, nil
}
