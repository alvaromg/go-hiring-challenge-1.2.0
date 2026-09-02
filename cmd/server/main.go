package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/lib/monitor"
	"github.com/mytheresa/go-hiring-challenge/models"
)

func main() {

	// create monitor

	logger, err := monitor.NewLogger("INFO", false)
	if err != nil {
		log.Fatalf("Error creating logger: %s", err)
	}

	monitor := monitor.NewMonitor(logger)

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_LOG_LEVEL"),
	)
	defer close()

	// Initialize repositories
	productsRepo := models.NewProductsRepository(db)
	categoriesRepo := models.NewCategoriesRepository(db)

	// Initialize applications
	catalogApp := catalog.NewCatalogApp(productsRepo, categoriesRepo)

	// Set up routing
	router := api.NewApiRouter(monitor, catalogApp)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: router,
	}

	// Start the server
	go func() {
		logger.Infof("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		logger.Infof("Server stopped gracefully")
	}()

	<-ctx.Done()
	logger.Infof("Shutting down server...")
	srv.Shutdown(ctx)
	stop()
}
