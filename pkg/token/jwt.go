package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(secret string, claims map[string]any) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	for k, v := range claims {
		tokenClaims[k] = v
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString, secret string) (map[string]any, error) {
	claims := map[string]any{}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	// get claims
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		claims = mapClaims
	}

	return claims, nil
}
