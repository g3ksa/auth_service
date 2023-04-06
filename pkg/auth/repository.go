package auth

import "AuthService/pkg/models"

type Repository interface {
	Refresh(userId int, refreshToken string) (string, error)
	Logout() error
	Get(email, password string) (*models.User, error)
}
