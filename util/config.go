package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Manually parse duration if it's a string like "15m"
	durationStr := viper.GetString("ACCESS_TOKEN_DURATION")
	config.AccessTokenDuration, err = time.ParseDuration(durationStr)
	if err != nil {
		return config, fmt.Errorf("invalid ACCESS_TOKEN_DURATION format: %w", err)
	}

	return config, nil
}
