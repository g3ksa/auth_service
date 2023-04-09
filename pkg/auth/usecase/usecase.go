package usecase

import (
	"AuthService/pkg/auth"
	"AuthService/pkg/hasher"
	"AuthService/pkg/parser"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
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

func (a *Authorizer) UpdateAccess(refreshToken string) (string, error) {
	claims, err := parser.ParseToken(refreshToken, []byte(a.refreshSecret))
	if err != nil {
		return "", err
	}

	if user, err := a.repo.Token(claims.UserId, refreshToken); err == nil {
		fmt.Println(claims.UserId)
		return jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: jwt.At(time.Now().Add(a.expireAccessDuration)),
				IssuedAt:  jwt.At(time.Now()),
			},
			UserId:      user.Id,
			UserCompany: user.UserCompany,
			UserRights:  user.UserRights,
			UserRole:    user.UserRole,
		}).SignedString([]byte(a.accessSecret))
	} else {
		return "", err
	}
}

func (a *Authorizer) SignOut(refreshToken string) error {
	claims, err := parser.ParseToken(refreshToken, []byte(a.refreshSecret))
	if err != nil {
		return err
	}
	return a.repo.Delete(claims.UserId)
}

func (a *Authorizer) SignIn(email, password string) ([]string, error) {

	user, err := a.repo.Get(auth.GetParams{
		Email:    email,
		Password: password,
	})

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
		UserRole:    user.UserRole,
	}).SignedString([]byte(a.accessSecret))

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: user.Id,
	}).SignedString([]byte(a.refreshSecret))

	hashedRt := hasher.HashString(refreshToken, []byte(a.refreshSecret))

	err = a.repo.PutToken(hashedRt, user.Id)

	return []string{accessToken, refreshToken}, err
}
