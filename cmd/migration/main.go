package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/canonflow/backend-starter/internal/app"
	_ "github.com/canonflow/backend-starter/internal/app"
	"github.com/canonflow/backend-starter/internal/config"
	_ "github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/core"
	"github.com/canonflow/backend-starter/pkg/helpers"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

var AVAILABLE_DRIVERS = []string{"mysql"}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	logger := core.GetLogger().With(
		zap.String("system", "migration"),
	)

	m := newMigrate(logger)
	defer m.Close()
	command := args[0]

	switch command {
	case "up":
		runUp(m, logger)
	case "down":
		if len(args) > 1 {
			steps, err := strconv.Atoi(args[1])
			if err != nil {
				logger.Fatal(fmt.Sprintf("Error: invalid steps '%s', must be an integer", args[1]))
			}
			runDownSteps(m, logger, steps)
		} else {
			runDown(m, logger)
		}
	case "version":
		runVersion(m, logger)
	case "force":
		if len(args) < 2 {
			logger.Fatal("Error: force requires a version argument. Usage: go run cmd/migrate/main.go force <version>")
		}
		version, err := strconv.Atoi(args[1])
		if err != nil {
			logger.Fatal(fmt.Sprintf("Error: invalid version '%s', must be an integer", args[1]))
		}
		runForce(m, logger, version)
	case "drop":
		runDrop(m, logger)
	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func newMigrate(logger *zap.Logger) *migrate.Migrate {
	driver := strings.ToLower(config.Get[string](config.DBDriver))
	username := config.Get[string](config.DBUsername)
	password := config.Get[string](config.DBPassword)
	host := config.Get[string](config.DBHost)
	port := config.Get[string](config.DBPort)
	database := config.Get[string](config.DBName)
	idleConnection := config.Get[int](config.DBIdle)
	maxConnection := config.Get[int](config.DBMax)
	maxLifetimeConnection := config.Get[int](config.DBLifeTime)

	if !helpers.SliceContains(AVAILABLE_DRIVERS, driver) {
		logger.Fatal("Unsupported Database Driver")
	}

	dsn := app.GetDriverConnection(app.DBProperty{
		Driver:                driver,
		Username:              username,
		Password:              password,
		Database:              database,
		Host:                  host,
		Port:                  port,
		IdleConnection:        idleConnection,
		MaxConnection:         maxConnection,
		MaxLifetimeConnection: maxLifetimeConnection,
	})

	migrationPath := fmt.Sprintf("file://./migrations/%s", driver)

	m, err := migrate.New(migrationPath, fmt.Sprintf("%s://%s", driver, dsn))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating migrate instance: %v", err))
	}

	return m
}

func runUp(m *migrate.Migrate, logger *zap.Logger) {
	logger.Info("Running migrations up...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(fmt.Sprintf("Error running up: %v", err))
	}
	logger.Info("Migrations up completed successfully.")
}

func runDown(m *migrate.Migrate, logger *zap.Logger) {
	logger.Info("Running migrations down...")
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(fmt.Sprintf("Error running down: %v", err))
	}
	logger.Info("Migrations down completed successfully.")
}

func runDownSteps(m *migrate.Migrate, logger *zap.Logger, steps int) {
	logger.Info(fmt.Sprintf("Rolling back %d migration(s)...", steps))
	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(fmt.Sprintf("Error rolling back %d steps: %v", steps, err))
	}
	logger.Info(fmt.Sprintf("Rolled back %d migration(s) successfully.", steps))
}

func runVersion(m *migrate.Migrate, logger *zap.Logger) {
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logger.Fatal(fmt.Sprintf("Error getting version: %v", err))
	}
	if err == migrate.ErrNilVersion {
		logger.Info("No migrations have been run yet.")
		return
	}
	logger.Info(fmt.Sprintf("Current version: %d | Dirty: %v", version, dirty))
}

func runForce(m *migrate.Migrate, logger *zap.Logger, version int) {
	logger.Info(fmt.Sprintf("Forcing migration to version %d...", version))
	if err := m.Force(version); err != nil {
		logger.Fatal(fmt.Sprintf("Error forcing version: %v", err))
	}
	logger.Info(fmt.Sprintf("Forced to version %d successfully.", version))
}

func runDrop(m *migrate.Migrate, logger *zap.Logger) {
	logger.Info("Dropping all migrations...")
	if err := m.Drop(); err != nil {
		logger.Fatal(fmt.Sprintf("Error dropping migrations: %v", err))
	}
	logger.Info("All migrations dropped successfully.")
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up               Apply all pending migrations")
	fmt.Println("  down [steps]     Revert all or N applied migrations (e.g. down 2)")
	fmt.Println("  version          Show current migration version")
	fmt.Println("  force <version>  Force set migration version (use after dirty state)")
	fmt.Println("  drop             Drop all migrations")
}
