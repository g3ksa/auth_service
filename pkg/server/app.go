package server

import (
	"AuthService/pkg/auth"
	"AuthService/pkg/auth/delivery"
	"AuthService/pkg/auth/repository/mysql"
	"AuthService/pkg/auth/usecase"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer  *http.Server
	authUseCase auth.UseCase
}

func NewApp() *App {
	db, err := sql.Open("mysql", "username:password@tcp(host:port)/db")
	if err != nil {
		log.Fatal(err)
	}
	userRepo := mysql.New(db)

	accessSecret, _ := os.LookupEnv("accessSecret")
	refreshSecret, _ := os.LookupEnv("refreshSecret")

	authUseCase := usecase.New(
		userRepo,
		accessSecret,
		refreshSecret,
		time.Minute,
	)
	return &App{
		authUseCase: authUseCase,
	}
}

func (a *App) Run(port string) error {
	router := gin.Default()

	router.Use(gin.Recovery(), gin.Logger())

	api := router.Group("/auth")
	delivery.RegisterHTTPEndpoints(api, a.authUseCase)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
