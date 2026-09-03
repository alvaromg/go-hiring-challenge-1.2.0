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
	"gorm.io/plugin/dbresolver"
)

// Config holds the connection parameters for a single postgres instance.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c Config) dsn() string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, sslMode)
}

// New opens a database connection that routes writes to write and reads to
// read via gorm's dbresolver plugin.
func New(write, read Config, pool PoolConfig, logLevel string, tracer trace.Tracer) (db *gorm.DB, close func() error) {
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

	db, err = gorm.Open(postgres.Open(write.dsn()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	resolver := dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{postgres.Open(write.dsn())},
		Replicas: []gorm.Dialector{postgres.Open(read.dsn())},
		Policy:   dbresolver.RandomPolicy{},
	})
	resolver.
		SetMaxOpenConns(pool.MaxOpenConns).
		SetMaxIdleConns(pool.MaxIdleConns).
		SetConnMaxLifetime(pool.ConnMaxLifetime).
		SetConnMaxIdleTime(pool.ConnMaxIdleTime)

	err = db.Use(resolver)
	if err != nil {
		log.Fatalf("failed to register db resolver: %s", err)
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
