package delivery

const (
	STATUS_OK    = "ok"
	STATUS_ERROR = "error"
)

type handler struct {
	useCase auth.UseCase
}
