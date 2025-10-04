package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string         `env:"ENVIRONMENT" env-default:"dev"`
	HTTP        HTTPConfig     `env-prefix:"HTTP_"`
	Postgres    PostgresConfig `env-prefix:"PG_"`
}

type HTTPConfig struct {
	Addr           string        `env:"ADDR"          env-default:":8080"`
	ReadTimeout    time.Duration `env:"READ_TIMEOUT"  env-default:"15s"`
	WriteTimeout   time.Duration `env:"WRITE_TIMEOUT" env-default:"15s"`
	IdleTimeout    time.Duration `env:"IDLE_TIMEOUT"  env-default:"120s"`
	HandlerTimeout time.Duration `env:"HANDLER_TIMEOUT" env-default:"2s"`
}
type PostgresConfig struct {
	Host            string        `env:"HOST"              env-default:"localhost"`
	Port            int           `env:"PORT"              env-default:"5432"`
	User            string        `env:"USER"              env-default:"postgres"`
	Password        string        `env:"PASSWORD"          env-default:"postgres"`
	DBName          string        `env:"DBNAME"            env-default:"spycat"`
	SSLMode         string        `env:"SSLMODE"           env-default:"disable"`
	ConnTimeout     time.Duration `env:"CONN_TIMEOUT"      env-default:"5s"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS"    env-default:"10"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"    env-default:"10"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" env-default:"30m"`
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode,
	)
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	return &cfg, nil
}
