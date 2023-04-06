package delivery

import "AuthService/pkg/auth"

const (
	STATUS_OK    = "ok"
	STATUS_ERROR = "error"
)

type handler struct {
	useCase auth.UseCase
}
