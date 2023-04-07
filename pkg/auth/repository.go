package auth

import (
	"AuthService/pkg/models"
)

type GetParams struct {
	Email    string
	Password string
}

type Repository interface {
	Refresh(userId int, refreshToken string) (*models.User, error)
	Logout(int) error
	Get(params GetParams) (*models.User, error)
}
