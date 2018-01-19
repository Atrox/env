package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	appEnvironment = "APP_ENV"
	devEnvironment = "development"
)

func init() {
	// error is ignored because we don't care if the .env is found and loaded,
	// in production there isn't one anyway
	godotenv.Load()
}

// Get returns the value of the specified environment variable
func Get(key string) string {
	return os.Getenv(key)
}

// GetDefault returns the value of the specified environment variable
// If the variable is not defined or is empty, return the default value
func GetDefault(key, defaultValue string) string {
	value := Get(key)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

// AppEnvironment returns the name of the current environment the app is in
func AppEnvironment() string {
	return strings.ToLower(Get(appEnvironment))
}

// IsDevelopment returns true if the current environment is development
func IsDevelopment() bool {
	return AppEnvironment() == devEnvironment
}

// IsProduction returns if the current environment is production
func IsProduction() bool {
	return !IsDevelopment()
}
