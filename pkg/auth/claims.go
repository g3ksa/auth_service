package auth

import (
	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	UserId      int      `json:"user_id"`
	UserCompany string   `json:"user_company"`
	UserRights  []string `json:"user_rights"`
	UserRole    string   `json:"user_role"`
}
