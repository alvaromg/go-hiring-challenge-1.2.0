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

	// Bootstrap OpenTelemetry: traces and logs are exported via OTLP/HTTP to
	// the grafana-lgtm container started by docker-compose.
	tracer, loggerProvider, shutdownOTel, err := monitor.SetupOTelSDK(ctx, serviceName, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatalf("Error setting up OpenTelemetry: %s", err)
	}
	defer func() {
		if err := shutdownOTel(context.Background()); err != nil {
			log.Printf("Error shutting down OpenTelemetry: %s", err)
		}
	}()

	// create monitor
	logger, err := monitor.NewLogger("INFO", false, loggerProvider)
	if err != nil {
		log.Fatalf("Error creating logger: %s", err)
	}

	monitor := monitor.NewMonitor(logger, tracer)

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_LOG_LEVEL"),
		tracer,
	)
	//nolint: errcheck
	defer close()

	// Initialize repositories
	productsRepo := models.NewProductsRepository(db)
	categoriesRepo := models.NewCategoriesRepository(db)

	// Initialize application
	catalogApp := catalog.NewCatalogApp(productsRepo, categoriesRepo)

	// Set up routing
	router := api.NewApiRouter(monitor, catalogApp)

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
