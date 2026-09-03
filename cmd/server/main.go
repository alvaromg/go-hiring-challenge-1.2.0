package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/infra/api"
	"github.com/mytheresa/go-hiring-challenge/infra/database"
	"github.com/mytheresa/go-hiring-challenge/infra/models"
	"github.com/mytheresa/go-hiring-challenge/infra/monitor"
	"github.com/mytheresa/go-hiring-challenge/infra/rest"
	_ "github.com/mytheresa/go-hiring-challenge/swagger"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Go Hiring Challenge
// @version 1.0
// @description This is API server for Go Hiring Challenge
// @termsOfService http://swagger.io/terms/

// @host localhost:8484
// @BasePath /
const serviceName = "go-hiring-challenge"

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

	runDocServer, err := strconv.ParseBool(os.Getenv("RUN_DOC_SERVER"))
	if err != nil {
		log.Fatalf("Unable to parse START_DOC_SERVER config value: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	writeConfig := database.Config{
		Host:     os.Getenv("POSTGRES_WRITE_HOST"),
		Port:     os.Getenv("POSTGRES_WRITE_PORT"),
		User:     os.Getenv("POSTGRES_WRITE_USER"),
		Password: os.Getenv("POSTGRES_WRITE_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_WRITE_DB"),
	}
	readConfig := database.Config{
		Host:     os.Getenv("POSTGRES_READ_HOST"),
		Port:     os.Getenv("POSTGRES_READ_PORT"),
		User:     os.Getenv("POSTGRES_READ_USER"),
		Password: os.Getenv("POSTGRES_READ_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_READ_DB"),
	}

	db, close := database.New(writeConfig, readConfig, database.PoolConfigFromEnv(), os.Getenv("POSTGRES_LOG_LEVEL"))
	//nolint: errcheck
	defer close()

	// Initialize repositories
	productsRepo := models.NewProductsRepository(db)
	categoriesRepo := models.NewCategoriesRepository(db)

	// Initialize application
	catalogApp := catalog.NewCatalogApp(logger, productsRepo, categoriesRepo)

	// Set up routing
	corsAllowedOrigins := rest.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	router := api.NewApiRouter(monitor, catalogApp, corsAllowedOrigins)

	// Set up the API server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: router,
	}

	// Start API server
	go func() {
		logger.Infof("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		logger.Infof("Server stopped gracefully")
	}()

	// Set up Swagger documentation server

	docRouter := http.NewServeMux()
	docRouter.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"),
	))
	docSrv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", "1323"),
		Handler: docRouter,
	}

	if runDocServer {
		// Start DOC server (swagger)
		go func() {
			logger.Infof("Starting Doc server on http://%s", docSrv.Addr)
			if err := docSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Swagger server failed: %s", err)
			}

			logger.Infof("Doc server stopped gracefully")
		}()
	}

	// Shutdown servers

	<-ctx.Done()
	logger.Infof("Shutting down API server...")
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("API server shutdown failed: %s", err)
	}
	if runDocServer {
		logger.Infof("Shutting down Doc server...")
		err = docSrv.Shutdown(ctx)
		if err != nil {
			log.Fatalf("Doc server shutdown failed: %s", err)
		}
	}

	stop()
}
