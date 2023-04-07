package app

import (
	"fmt"

	env "github.com/Netflix/go-env"
)

type Config struct {
	DBHost      string `env:"DB_HOST"`
	DBUser      string `env:"DB_USER"`
	DBPass      string `env:"DB_PASS"`
	DBPort      string `env:"DB_PORT"`
	DBName      string `env:"DB_NAME"`
	HTTPAddress string `env:"HTTP_ADDRESS"`
}

func NewConfig() (Config, error) {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (c Config) migrateDatabaseDSN() string {
	return c.databaseConn("pgx5")
}

func (c Config) ConnectDatabaseDSN() string {
	return c.databaseConn("postgresql")
}

func (c Config) databaseConn(driver string) string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		driver,
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
