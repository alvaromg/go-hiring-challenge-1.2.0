package database

import (
	"log"
	"os"
	"strconv"
	"time"
)

// PoolConfig holds connection pool tuning options. It applies to both the
// write and read connections managed by dbresolver, so it's kept separate
// from Config rather than duplicated on write/read. A zero value leaves the
// corresponding setting at the database/sql driver default (unlimited).
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// PoolConfigFromEnv reads pool tuning options from environment. An unset or empty variable keeps the
// corresponding PoolConfig field at its zero value (driver default).
func PoolConfigFromEnv() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    envInt("POSTGRES_MAX_OPEN_CONNS"),
		MaxIdleConns:    envInt("POSTGRES_MAX_IDLE_CONNS"),
		ConnMaxLifetime: envDuration("POSTGRES_CONN_MAX_LIFETIME"),
		ConnMaxIdleTime: envDuration("POSTGRES_CONN_MAX_IDLE_TIME"),
	}
}

func envInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("failed to parse %s as int: %s", key, err)
	}

	return n
}

func envDuration(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("failed to parse %s as duration: %s", key, err)
	}

	return d
}
