package config

import (
	"os"
	"strconv"
)

// Config provides app configuration accessors
type Config struct{}

func New() *Config {
	return &Config{}
}

func (c *Config) Get(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func (c *Config) GetBool(key string, defaultValue bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return defaultValue
}

func (c *Config) GetInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}
