package main

import (
	"boilerplate/internal/container"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second

	successCode = iota - 1
	errorCode
)

func main() {
	errorCode := errorCode
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println("error on logger sync:", err)
		}
		os.Exit(errorCode)
	}()

	logger.Info("starting the application")
	app := container.NewApplication(ctx, logger)

	router := gin.Default()

	router.Use(ginzap.Ginzap(app.Log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(app.Log, true))
	router.Use(otelgin.Middleware("server"))

	router.GET("/healthcheck", app.HandleHeathCheck)
	router.GET("/readiness", app.HandleReadiness)
	router.POST("/echo", app.HandleEcho)

	srv := &http.Server{
		ReadHeaderTimeout: shutdownTimeout,
		Addr:              ":8080",
		Handler:           router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error on server:", zap.Error(err))
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	log.Println("shutting down the server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("error on server shutdown:", zap.Error(err))
		return
	}

	log.Println("shutting down the application")
	if err := app.GracefulShutdown(ctx); err != nil {
		logger.Error("error on application shutdown:", zap.Error(err))
		return
	}

	errorCode = successCode
}
