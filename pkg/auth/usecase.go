package auth

type UseCase interface {
	UpdateAccess(refreshToken string) (string, error)
	SignOut(refreshToken string) error
	SignIn(email, password string) ([]string, error)
}
