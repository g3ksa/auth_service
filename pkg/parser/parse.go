package parser

import (
	"AuthService/pkg/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"net/http"
)

func ParseToken(accessToken string, signingKey []byte) (string, *auth.Response) {
	token, err := jwt.ParseWithClaims(accessToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return "", &auth.Response{
			Error:   err,
			Status:  http.StatusUnauthorized,
			Message: "Неверный формат токена",
		}
	}

	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", &auth.Response{
		Error:   auth.ErrInvalidAccessToken,
		Status:  http.StatusUnauthorized,
		Message: "Неверный формат токена",
	}
}
