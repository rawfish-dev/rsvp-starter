package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	defaultHTTPPort = 6001
	sessionDuration = time.Minute * 20
)

// Config holds necessary config values.
type Config struct {
	HTTPPort int
	Postgres PostgresConfig
	Session  SessionConfig
	JWT      JWTConfig
}

// PostgresConfig contains the connection URL and other DB options.
type PostgresConfig struct {
	URL            string
	MaxIdle        int
	MaxConnections int
}

// SessionConfig contains the duration of each valid session.
type SessionConfig struct {
	Duration time.Duration
}

// JWTConfig contains the config values required to create valid JWTs.
type JWTConfig struct {
	HMACSecret  string
	TokenIssuer string
}

var (
	once   sync.Once
	config Config
)

// LoadConfig instantiates a singleton object that holds necessary config values.
// This function panics if required environment values are not set properly.
func LoadConfig() Config {
	once.Do(func() {
		config = Config{
			HTTPPort: parseHTTPPort(),
			Postgres: loadPostgresConfig(),
			Session:  loadSessionConfig(),
			JWT:      loadJWTConfig(),
		}
	})

	return config
}

func parseHTTPPort() int {
	httpPortStr, ok := os.LookupEnv("HTTP_PORT")
	if !ok || httpPortStr == "" {
		return defaultHTTPPort
	}

	httpPort, err := strconv.ParseInt(httpPortStr, 10, 32)
	if err != nil {
		logrus.Fatalf("HTTP_PORT value '%s' could not be parsed due to %s", httpPortStr, err.Error())
	}

	return int(httpPort)
}

func loadPostgresConfig() PostgresConfig {
	postgresURL, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		logrus.Fatal("POSTGRES_URL not set")
	}

	return PostgresConfig{
		URL:            postgresURL,
		MaxIdle:        15,
		MaxConnections: 15,
	}
}

func loadSessionConfig() SessionConfig {
	return SessionConfig{
		Duration: sessionDuration,
	}
}

func loadJWTConfig() JWTConfig {
	hmacSecret, ok := os.LookupEnv("HMAC_SECRET")
	if !ok {
		logrus.Fatal("HMAC_SECRET not set")
	}

	tokenIssuer, ok := os.LookupEnv("TOKEN_ISSUER")
	if !ok {
		logrus.Fatal("TOKEN_ISSUER not set")
	}

	return JWTConfig{
		HMACSecret:  hmacSecret,
		TokenIssuer: tokenIssuer,
	}
}
