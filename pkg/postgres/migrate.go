package postgres

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
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
			log.Fatalf("Migrate: postgres connect error: %s", err)
		}

		err = m.Up()
		defer m.Close()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migrate: up error: %s", err)
		}

		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Migrate: no change")
		} else {
			log.Printf("Migrate: up success")
		}

		return clickhouseMigrate(config)
	}
}
