package usecase

import (
	"AuthService/pkg/auth"
	"github.com/dgrijalva/jwt-go/v4"
	_ "golang.org/x/crypto/bcrypt"
	"time"
)

type Authorizer struct {
	repo auth.Repository

	accessSecret         string
	refreshSecret        string
	expireAccessDuration time.Duration
}

func New(repo auth.Repository, accessSecret, refreshSecret string, expireAccessDuration time.Duration) *Authorizer {
	return &Authorizer{
		repo:                 repo,
		accessSecret:         accessSecret,
		refreshSecret:        refreshSecret,
		expireAccessDuration: expireAccessDuration,
	}
}

func (a *Authorizer) SignIn(email, password string) ([]string, error) {
	user, err := a.repo.Get(email, password)

	if err != nil {
		return []string{}, err
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireAccessDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId:      user.Id,
		UserCompany: user.UserCompany,
		UserRights:  user.UserRights,
	}).SignedString([]byte(a.accessSecret))

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: user.Id,
	}).SignedString([]byte(a.refreshSecret))

	return []string{accessToken, refreshToken}, err
}
