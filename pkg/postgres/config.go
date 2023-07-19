package postgres

import (
	"fmt"
	"time"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Config хранит данные для подключения постгреса
type Config struct {
	User          string        `yaml:"user" env:"PG_USER"`
	Password      string        `yaml:"password" env:"PG_PASSWORD"`
	Host          string        `yaml:"host" env:"PG_HOST"`
	Port          string        `yaml:"port" env:"" envDefault:"5432"`
	DBName        string        `yaml:"db_name" env:"PG_DB_NAME"`
	MaxPoolSize   int           `yaml:"max_pool_size" env:"pg_max_pool_size"`
	ConnAttempts  int           `yaml:"conn_attempts" env:"pg_ConnAttempts "`
	ConnTimeout   time.Duration `yaml:"conn_timeout" env:"pg_ConnTimeout  "`
	ClickhouseURL string        `yaml:"clickhouse_url" env:"CLICKHOUSE_URL"`
}

// NewConfig конструтор
func NewConfig() Config {
	return Config{
		Host:         "localhost",
		Port:         "5432",
		MaxPoolSize:  _defaultMaxPoolSize,
		ConnAttempts: _defaultConnAttempts,
		ConnTimeout:  _defaultConnTimeout,
	}
}

// URL возвращает конфиг в виде строки
func (c *Config) URL() string {
	// postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}
