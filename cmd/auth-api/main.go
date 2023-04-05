package main

import (
	"AuthService/config"
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

}
