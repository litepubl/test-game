package postgres

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
)

func clickhouseMigrate(config *Config) error {
	databaseURL := config.ClickhouseURL

	var (
		attempts = config.ConnAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations/ch", databaseURL)
		if err == nil {
			break
		}

		log.Printf("clickhouseMigrate: clickhouse is trying to connect, attempts left: %d", attempts)
		time.Sleep(config.ConnTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("clickhouseMigrate: clickhouse connect error: %s", err)
	}

	err = m.Up()
	defer m.Close() //nolint
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("clickhouseMigrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("clickhouseMigrate: no change")
		return nil
	}

	log.Printf("clickhouseMigrate: up success")
	return nil
}
