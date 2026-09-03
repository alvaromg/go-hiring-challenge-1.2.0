package database

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New opens the database connection and registers a tracing plugin that
// starts an OpenTelemetry span for every SQL statement GORM sends. Pass nil
// for tracer where tracing isn't needed (e.g. tests, one-off scripts); a
// no-op tracer is used in that case.
func New(user, password, dbname, port, logLevel string, tracer trace.Tracer) (db *gorm.DB, close func() error) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	dbLogLevel, err := parseLogLevel(logLevel)
	if err != nil {
		log.Fatalf("failed to parse database log level: %s", err)
	}

	newLogger := logger.New(
		newJSONLogWriter(),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  dbLogLevel,  // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %s", err)
	}

	if tracer == nil {
		tracer = noop.NewTracerProvider().Tracer("database")
	}
	if err := db.Use(newTracingPlugin(tracer)); err != nil {
		log.Fatalf("failed to register database tracing plugin: %s", err)
	}

	return db, sqlDB.Close
}
