package delivery

import (
	"AuthService/pkg/auth"
	"AuthService/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	STATUS_OK    = "ok"
	STATUS_ERROR = "error"
)

type response struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func newResponse(status, msg string) *response {
	return &response{
		Status: status,
		Msg:    msg,
	}
}

type handler struct {
	useCase auth.UseCase
}

func newHandler(usecase auth.UseCase) *handler {
	return &handler{
		useCase: usecase,
	}
}

type responseWithToken struct {
	*response
	Token string `json:"token,omitempty"`
}

func newResponseWithToken(status, msg, token string) *responseWithToken {
	return &responseWithToken{
		&response{
			Status: status,
			Msg:    msg,
		},
		token,
	}
}

func (h *handler) signIn(c *gin.Context) {
	inp := new(models.User)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//fmt.Println(inp.Email, inp.Password)
	tokens, err := h.useCase.SignIn(inp.Email, inp.Password)
	if err != nil {
		if err == auth.ErrIncorrectPassword || err == auth.ErrUserDoesNotExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, newResponseWithToken(STATUS_ERROR, "Неверные данные пользователя", ""))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponseWithToken(STATUS_ERROR, err.Error(), ""))
		return
	}
	c.SetCookie("refreshToken", tokens[1], 3600, "/", ".", false, true)
	c.JSON(http.StatusOK, newResponseWithToken(STATUS_OK, "", tokens[0]))
}

func (h *handler) signOut(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.useCase.SignOut(refreshToken)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (h *handler) token(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newAccessToken, err := h.useCase.UpdateAccess(refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, newResponseWithToken(STATUS_ERROR, "Не удалось авторизовать пользователя", ""))
		return
	}
	c.JSON(http.StatusOK, newResponseWithToken(STATUS_OK, "", newAccessToken))
}
