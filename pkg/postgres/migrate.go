package postgres

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	migrateFunc = func(config *Config) error {
		databaseURL := config.URL() + "?sslmode=disable"

		var (
			attempts = config.ConnAttempts
			err      error
			m        *migrate.Migrate
		)

		for attempts > 0 {
			m, err = migrate.New("file://migrations/pg", databaseURL)
			if err == nil {
				break
			}

			log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
			time.Sleep(config.ConnTimeout)
			attempts--
		}

		if err != nil {
			log.Fatalf("Migrate: postgres connect error: %w", err)
		}

		err = m.Up()
		defer m.Close()

		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Printf("Migrate: no change")
			} else {
				return fmt.Errorf("Migrate: up error: %w", err)
			}
		} else {
			log.Printf("Migrate: up success")
		}

		return clickhouseMigrate(config)
	}
}
