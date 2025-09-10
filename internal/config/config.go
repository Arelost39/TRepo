package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL		string
	Port			string
	SystemPeriod	int64
	Time			time.Time
}

func Load() (*Config, error) {
	c := &Config{}

	// БДшка
	c.DatabaseURL = os.Getenv("DATABASE_URL")
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL не задано")
	}

	// порт сервера
	c.DatabaseURL = os.Getenv("PORT")
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("PORT не задано")
	}

	// окно для next_takings
	// будем считать, что окно задано в .env
	SystemPeriodStr := os.Getenv("SYSTEM_PERIOD")
	if SystemPeriodStr == "" {
		return nil, fmt.Errorf("SYSTEM_PERIOD не задано")
	}
	minutes, err := strconv.ParseInt(SystemPeriodStr, 64, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid SYSTEM_PERIOD: %w", err)
	}

	c.SystemPeriod = minutes

	return c, nil
}