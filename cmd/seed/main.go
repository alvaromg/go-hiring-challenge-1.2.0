package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"

	"github.com/mytheresa/go-hiring-challenge/infra/database"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

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

	db, close := database.New(writeConfig, readConfig, database.PoolConfigFromEnv(), os.Getenv("POSTGRES_LOG_LEVEL"), nil)
	//nolint: errcheck
	defer close()

	dir := os.Getenv("POSTGRES_SQL_DIR")
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("reading directory failed: %v", err)
	}

	// Filter and sort .sql files
	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, file := range sqlFiles {
		path := filepath.Join(dir, file.Name())

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("reading file %s failed: %v", file.Name(), err)
		}

		sql := string(content)
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("executing %s failed: %v", file.Name(), err)
			return
		}

		log.Printf("Executed %s successfully\n", file.Name())
	}
}
