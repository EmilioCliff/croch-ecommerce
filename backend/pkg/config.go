package pkg

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	MYSQL_USER             string        `mapstructure:"MYSQL_USER"`
	MYSQL_PASSWORD         string        `mapstructure:"MYSQL_PASSWORD"`
	MYSQL_DB               string        `mapstructure:"MYSQL_DB"`
	DB_DSN                 string        `mapstructure:"DB_DSN"`
	MIGRATION_PATH         string        `mapstructure:"MIGRATION_PATH"`
	TOKEN_DURATION         time.Duration `mapstructure:"TOKEN_DURATION"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TOKEN_SYMMETRY_KEY     string        `mapstructure:"TOKEN_SYMMETRY_KEY"`
	PASSWORD_COST          int           `mapstructure:"PASSWORD_COST"`
}

// TOKEN_DURATION="15m"
// TOKEN_SYMMETRY_KEY="12345678901234567890123456789012"

// Loads app configuration from .env file.
func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	err := viper.Unmarshal(&config)

	return config, err
}
