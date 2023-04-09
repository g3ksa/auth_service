package delivery

import (
	"AuthService/pkg/auth"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, useCase auth.UseCase) {
	h := newHandler(useCase)

	router.POST("/login", h.signIn)
	router.POST("/refresh", h.token)
	router.POST("/logout", h.signOut)
}
