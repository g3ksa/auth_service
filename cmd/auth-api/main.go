package main

import (
	"AuthService/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	if err := config.Init(); err != nil {
		log.Fatalf(
			"%s",
			err.Error(),
		)
	}
}

func main() {
	router := gin.Default()

	router.GET("/cookie", func(c *gin.Context) {

		refreshToken, err := c.Cookie("refreshToken")

		if err != nil {
			refreshToken = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		fmt.Printf("Cookie value: %s \n", refreshToken)
	})

	router.Run()
}
