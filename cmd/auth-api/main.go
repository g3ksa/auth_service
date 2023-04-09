package main

import (
	"AuthService/pkg/server"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	app := server.NewApp()
	port, _ := os.LookupEnv("port")
	fmt.Println(port)
	if err := app.Run(port); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
