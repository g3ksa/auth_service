package auth

import "AuthService/pkg/models"

type Repository interface {
	GetToken(accessToken string) (string, error)
	Delete(userId int) error
	GetUser(email, password string) (*models.User, error)
}
