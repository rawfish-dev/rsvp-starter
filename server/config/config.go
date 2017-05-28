package config

import (
	"os"

	"github.com/Sirupsen/logrus"
)

type Config struct {
	DebugEnabled bool           `json:"debugEnabled"`
	HttpPort     int            `json:"httpPort"`
	Postgres     PostgresConfig `json:"postgres"`
	Session      SessionConfig  `json:"session"`
	JWT          JWTConfig      `json:"jwt"`
}

type PostgresConfig struct {
	URL            string `json:"url"`
	MaxIdle        int    `json:"maxIdle"`
	MaxConnections int    `json:"maxConnections"`
}

type SessionConfig struct {
	DurationSeconds int `json:"durationSeconds"`
}

type JWTConfig struct {
	HMACSecret  string `json:"hmacSecret"`
	TokenIssuer string `json:"tokenIssuer"`
}

var (
	config *Config

	validEnvironments = []string{"test", "development", "production"}
)

func Load() *Config {

	if config == nil {

		// Check the environment and load the corresponding file
		currentEnvironment := getEnvironment()

		logrus.Infof("config - loading configuration for %s", currentEnvironment)

		switch currentEnvironment {
		case "test":
			config = TestConfig()

		case "development":
			config = developmentConfig()

		case "production":
			config = productionConfig()

		default:
			config = developmentConfig()
		}
	}

	return config
}

func getEnvironment() string {

	currentEnvironment := os.Getenv("WEDDING_RSVP_ENV")

	for _, env := range validEnvironments {
		if currentEnvironment == env {
			return currentEnvironment
		}
	}

	// Return default
	return "development"
}

func TestConfig() *Config {
	logrus.Info("Loading test configuration...")

	return &Config{
		HttpPort: 6001,
		Postgres: PostgresConfig{
			URL:            "postgres://postgres@localhost/wedding_rsvp_test?sslmode=disable",
			MaxIdle:        15,
			MaxConnections: 15,
		},
		Session: SessionConfig{
			DurationSeconds: 600,
		},
		JWT: JWTConfig{
			HMACSecret:  "secret",
			TokenIssuer: "weddingRSVPTest",
		},
	}
}

func developmentConfig() *Config {
	logrus.Info("Loading development configuration...")

	// Temporarily use in code values
	return &Config{
		HttpPort: 5000,
		Postgres: PostgresConfig{
			URL:            "postgres://postgres@localhost/wedding_rsvp_development?sslmode=disable",
			MaxIdle:        15,
			MaxConnections: 15,
		},
		Session: SessionConfig{
			DurationSeconds: 3600,
		},
		JWT: JWTConfig{
			HMACSecret:  "reallysecret",
			TokenIssuer: "weddingRSVPDevelopment",
		},
	}
}

func productionConfig() *Config {
	logrus.Info("Loading production configuration...")

	productionDBURL, ok := os.LookupEnv("POSTGRES_URL")
	if !ok || productionDBURL == "" {
		panic("POSTGRES_URL not set")
	}

	hmacSecret, ok := os.LookupEnv("HMAC_SECRET")
	if !ok || hmacSecret == "" {
		panic("HMAC_SECRET not set")
	}

	return &Config{
		HttpPort: 5000,
		Postgres: PostgresConfig{
			URL:            productionDBURL,
			MaxIdle:        15,
			MaxConnections: 15,
		},
		Session: SessionConfig{
			DurationSeconds: 3600,
		},
		JWT: JWTConfig{
			HMACSecret:  hmacSecret,
			TokenIssuer: "weddingRSVPProduction",
		},
	}
}
