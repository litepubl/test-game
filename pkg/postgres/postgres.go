package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Postgres  агрегатор пула коннектов и построителя запросов
type Postgres struct {
	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

var migrateFunc func(config *Config) error

// New конструтор
func New(config *Config) (*Postgres, error) {
	pg := &Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	if migrateFunc != nil {
		err := migrateFunc(config)
		if err != nil {
			return nil, err
		}
	}

	poolConfig, err := pgxpool.ParseConfig(config.URL())
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(config.MaxPoolSize)
	connAttempts := config.ConnAttempts
	for connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", connAttempts)
		time.Sleep(config.ConnTimeout)
		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close закрывает соеденение с БД
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
