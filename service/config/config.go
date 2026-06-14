package config

import (
	"fmt"
	"strconv"
	"time"
)

const (
	ENV_PROD    = "prod"
	ENV_PREPROD = "preprod"
	ENV_RECETTE = "recette"
)

type ServerConfig struct {
	ServiceName        string
	Environment        string
	PGDriver           string
	MSSQLDriver        string
	HTTPAddr           string
	HTTPReadTimeout    time.Duration
	HTTPWriteTimeout   time.Duration
	HTTPIdleTimeout    time.Duration
	ServerShutdownTime int
	Cache              CacheConfig
}

// CacheConfig holds optional Valkey cache settings.
// When URL is empty the cache layer is disabled.
type CacheConfig struct {
	// URL is the Valkey connection string (valkey:// or valkeys:// for TLS).
	URL string
	// TLSCACert is an optional path to a PEM-encoded CA certificate used to
	// verify the server's TLS certificate (required for self-signed certs).
	TLSCACert string
	// Prefix is prepended to every cache key to namespace entries.
	// Defaults to the service name.
	Prefix string
	// TTL is the default cache TTL applied when no method-specific TTL is set.
	TTL time.Duration
	// HealthTTL overrides TTL for GetHealth cache entries.
	HealthTTL time.Duration
}

func LoadServerConfig(getenv func(string) string) (ServerConfig, error) {
	if getenv == nil {
		getenv = func(string) string { return "" }
	}

	env := map[string]string{
		"SERVICE_NAME":         getenv("SERVICE_NAME"),
		"ENVIRONMENT":          getenv("ENVIRONMENT"),
		"PG_DRIVER":            getenv("PG_DRIVER"),
		"MSSQL_DRIVER":         getenv("MSSQL_DRIVER"),
		"DATABASE_URL":         getenv("DATABASE_URL"),
		"MSSQL_DATABASE_URL":   getenv("MSSQL_DATABASE_URL"),
		"HTTP_ADDR":            getenv("HTTP_ADDR"),
		"HTTP_READ_TIMEOUT":    getenv("HTTP_READ_TIMEOUT"),
		"HTTP_WRITE_TIMEOUT":   getenv("HTTP_WRITE_TIMEOUT"),
		"HTTP_IDLE_TIMEOUT":    getenv("HTTP_IDLE_TIMEOUT"),
		"SERVER_SHUTDOWN_TIME": getenv("SERVER_SHUTDOWN_TIME"),
		"CACHE_URL":            getenv("CACHE_URL"),
		"CACHE_TLS_CA_CERT":    getenv("CACHE_TLS_CA_CERT"),
		"CACHE_PREFIX":         getenv("CACHE_PREFIX"),
		"CACHE_TTL":            getenv("CACHE_TTL"),
		"CACHE_HEALTH_TTL":     getenv("CACHE_HEALTH_TTL"),
	}

	pgDriver := valueOrDefault(env["PG_DRIVER"], "inmemory")
	mssqlDriver := valueOrDefault(env["MSSQL_DRIVER"], "inmemory")
	serviceName := valueOrDefault(env["SERVICE_NAME"], SERVICE_NAME)
	environment := valueOrDefault(env["ENVIRONMENT"], "recette")
	httpAddr := valueOrDefault(env["HTTP_ADDR"], ":8080")
	httpReadTimeout, err := durationOrDefault(env["HTTP_READ_TIMEOUT"], 5*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid HTTP_READ_TIMEOUT: %w", err)
	}
	httpWriteTimeout, err := durationOrDefault(env["HTTP_WRITE_TIMEOUT"], 5*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT: %w", err)
	}
	httpIdleTimeout, err := durationOrDefault(env["HTTP_IDLE_TIMEOUT"], 5*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid HTTP_IDLE_TIMEOUT: %w", err)
	}
	serverShutdownInSec, err := intOrDefault(env["SERVER_SHUTDOWN_TIME"], "5")
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid SERVER_SHUTDOWN_TIME: %w", err)
	}

	if err := ValidatePersistenceEnv(pgDriver, mssqlDriver, env); err != nil {
		return ServerConfig{}, err
	}
	if err := ValidateEnvironment(environment); err != nil {
		return ServerConfig{}, err
	}

	cacheCfg, err := loadCacheConfig(serviceName, env)
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		ServiceName:        serviceName,
		Environment:        environment,
		PGDriver:           pgDriver,
		MSSQLDriver:        mssqlDriver,
		HTTPAddr:           httpAddr,
		HTTPReadTimeout:    httpReadTimeout,
		HTTPWriteTimeout:   httpWriteTimeout,
		HTTPIdleTimeout:    httpIdleTimeout,
		ServerShutdownTime: serverShutdownInSec,
		Cache:              cacheCfg,
	}, nil
}

func loadCacheConfig(serviceName string, env map[string]string) (CacheConfig, error) {
	if env["CACHE_URL"] == "" {
		return CacheConfig{}, nil
	}

	cacheTTL, err := durationOrDefault(env["CACHE_TTL"], 60*time.Second)
	if err != nil {
		return CacheConfig{}, fmt.Errorf("invalid CACHE_TTL: %w", err)
	}

	cacheHealthTTL, err := durationOrDefault(env["CACHE_HEALTH_TTL"], cacheTTL)
	if err != nil {
		return CacheConfig{}, fmt.Errorf("invalid CACHE_HEALTH_TTL: %w", err)
	}

	return CacheConfig{
		URL:       env["CACHE_URL"],
		TLSCACert: env["CACHE_TLS_CA_CERT"],
		Prefix:    valueOrDefault(env["CACHE_PREFIX"], serviceName),
		TTL:       cacheTTL,
		HealthTTL: cacheHealthTTL,
	}, nil
}

func ValidateEnvironment(environment string) error {
	switch environment {
	case ENV_PROD, ENV_PREPROD, ENV_RECETTE:
		return nil
	default:
		return fmt.Errorf("unsupported ENVIRONMENT %q", environment)
	}
}

func ValidatePersistenceEnv(pgDriver, mssqlDriver string, env map[string]string) error {
	switch pgDriver {
	case "inmemory", "postgres":
	default:
		return fmt.Errorf("unsupported PG_DRIVER %q", pgDriver)
	}
	if pgDriver == "postgres" && env["DATABASE_URL"] == "" {
		return fmt.Errorf("missing DATABASE_URL for postgres driver")
	}

	switch mssqlDriver {
	case "inmemory", "mssql":
	default:
		return fmt.Errorf("unsupported MSSQL_DRIVER %q", mssqlDriver)
	}
	if mssqlDriver == "mssql" && env["MSSQL_DATABASE_URL"] == "" {
		return fmt.Errorf("missing MSSQL_DATABASE_URL for mssql driver")
	}

	return nil
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func intOrDefault(value, fallback string) (int, error) {
	if value == "" {
		return strconv.Atoi(fallback)
	}
	return strconv.Atoi(value)
}

func durationOrDefault(value string, fallback time.Duration) (time.Duration, error) {
	if value == "" {
		return fallback, nil
	}
	return time.ParseDuration(value)
}
