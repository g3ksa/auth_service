package parser

import (
	"AuthService/pkg/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
)

func ParseToken(accessToken string, signingKey []byte) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод шифрования: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке токена: %s", err)
	}

	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("что-то пошло не так")
}
