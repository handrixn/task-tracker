package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/handrixn/task-tracker/config"
	"github.com/handrixn/task-tracker/internal/handler"
	"github.com/handrixn/task-tracker/internal/repository"
	"github.com/handrixn/task-tracker/internal/router"
	"github.com/handrixn/task-tracker/internal/service"
	"github.com/spf13/viper"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Failed load env")
	}

	db, err := config.InitDB()

	if err != nil {
		log.Fatal("Error init Database connection")
	}

	taskRepository := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepository)
	taskHandler := handler.NewTaskHandler(taskService)
	mainRouter := router.NewRouter()
	router.NewTaskRouter(mainRouter, taskHandler)

	server := &http.Server{
		Addr:    fmt.Sprint(":", viper.GetString("APP_PORT")),
		Handler: mainRouter,
	}

	// Start the HTTP server in a separate goroutine
	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Handle graceful shutdown
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Block until a signal is received
	<-stop

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully shut down")
}
