package helper

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/mytheresa/go-hiring-challenge/infra/api"
	"github.com/mytheresa/go-hiring-challenge/infra/database"
	"github.com/mytheresa/go-hiring-challenge/infra/models"
	"github.com/mytheresa/go-hiring-challenge/infra/monitor"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"gorm.io/gorm"

	catalogapp "github.com/mytheresa/go-hiring-challenge/app/catalog"
)

const (
	dbUser     = "postgres"
	dbPassword = "password"
	dbName     = "challenge"
	sqlDir     = "../../sql"
)

func BuildApi() (http.Handler, func()) {
	queries, err := sortedSQLQueries(sqlDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read sql files: %s\n", err)
		os.Exit(1)
	}

	// dbUser/dbPassword match the preset's own defaults, so no WithUser
	// call is needed here (it would try to create a duplicate "postgres"
	// superuser and fail).
	opts := []postgres.Option{
		postgres.WithDatabase(dbName),
		postgres.WithQueries(queries...),
	}

	container, err := gnomock.Start(postgres.Preset(opts...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres container: %s\n", err)
		os.Exit(1)
	}

	config := database.Config{
		Host:     "localhost",
		Port:     fmt.Sprintf("%d", container.DefaultPort()),
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
	}

	// The test container is a single instance, so reads and writes share the same config.
	db, closeDB := database.New(config, config, database.PoolConfig{}, "error")

	router := newRouter(db)

	closeFn := func() {
		//nolint: errcheck
		closeDB()
		//nolint: errcheck
		gnomock.Stop(container)
	}

	return router, closeFn
}

// sortedSQLQueries reads the migration files under dir sorted lexically, so
// they get applied in the same numeric order they're meant to run in
// (000-truncate.sql, 001-products.sql, ...), and returns their contents in
// that order.
func sortedSQLQueries(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)

	queries := make([]string, len(files))
	for i, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		queries[i] = string(content)
	}

	return queries, nil
}

func newRouter(db *gorm.DB) http.Handler {
	productsRepo := models.NewProductsRepository(db)
	categoriesRepo := models.NewCategoriesRepository(db)
	noopMonitor := monitor.NewNoopMonitor()
	catalogApp := catalogapp.NewCatalogApp(noopMonitor.Logger(), productsRepo, categoriesRepo)

	return api.NewApiRouter(noopMonitor, catalogApp)
}
