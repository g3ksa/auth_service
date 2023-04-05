package usecase

import (
	"AuthService/pkg/auth"
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
