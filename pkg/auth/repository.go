package auth

import "AuthService/pkg/models"

type GetParams struct {
	Email    string
	Password string
}

type Repository interface {
	Token(userId int, refreshToken string) (*models.User, error)
	Delete(int) error
	Get(params GetParams) (*models.User, error)
	PutToken(refreshToken string, userId int) error
}
