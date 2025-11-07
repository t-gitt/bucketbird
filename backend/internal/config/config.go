package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	AppName           string
	Env               string
	HTTPPort          string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	AllowedOrigins    []string
	DBDSN             string
	S3Endpoint        string
	S3Region          string
	S3AccessKey       string
	S3SecretKey       string
	S3UseSSL          bool
	JWTSecret         string
	EncryptionKey     []byte
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
	CookieSecure      bool
	AllowRegistration bool
	EnableDemoLogin   bool
}

const (
	defaultAppName         = "bucketbird-api"
	defaultEnv             = "development"
	defaultHTTPPort        = "8080"
	defaultReadTimeout     = 30 * time.Minute // Allow reading large file uploads
	defaultWriteTimeout    = 30 * time.Minute // Allow writing large file uploads to S3
	defaultJWTSecret       = "bucketbird-insecure-dev-secret"
	defaultEncryptionKey   = "bucketbird-dev-key-32-bytes-long!!" // Must be exactly 32 bytes for AES-256
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 7 * 24 * time.Hour

	defaultDBHost     = "postgres"
	defaultDBPort     = "5432"
	defaultDBName     = "bucketbird"
	defaultDBUser     = "bucketbird"
	defaultDBPassword = "bucketbird"

	defaultS3Endpoint  = "http://minio:9000"
	defaultS3Region    = "us-east-1"
	defaultS3AccessKey = "minioadmin"
	defaultS3SecretKey = "minioadmin"
	defaultS3UseSSL    = false
)

func Load() Config {
	encKey := getEnv("BB_ENCRYPTION_KEY", defaultEncryptionKey)
	if len(encKey) != 32 {
		panic("BB_ENCRYPTION_KEY must be exactly 32 bytes for AES-256")
	}

	cfg := Config{
		AppName:           getEnv("BB_APP_NAME", defaultAppName),
		Env:               getEnv("BB_ENV", defaultEnv),
		HTTPPort:          getEnv("BB_HTTP_PORT", defaultHTTPPort),
		ReadTimeout:       getDurationEnv("BB_HTTP_READ_TIMEOUT", defaultReadTimeout),
		WriteTimeout:      getDurationEnv("BB_HTTP_WRITE_TIMEOUT", defaultWriteTimeout),
		AllowedOrigins:    []string{"*"},
		JWTSecret:         getEnv("BB_JWT_SECRET", defaultJWTSecret),
		EncryptionKey:     []byte(encKey),
		AccessTokenTTL:    getDurationEnv("BB_ACCESS_TOKEN_TTL", defaultAccessTokenTTL),
		RefreshTokenTTL:   getDurationEnv("BB_REFRESH_TOKEN_TTL", defaultRefreshTokenTTL),
		CookieSecure:      getBoolEnv("BB_COOKIE_SECURE", false),
		AllowRegistration: getBoolEnv("BB_ALLOW_REGISTRATION", true),
		EnableDemoLogin:   getBoolEnv("BB_ENABLE_DEMO_LOGIN", false),
	}

	if origins := strings.TrimSpace(os.Getenv("BB_ALLOWED_ORIGINS")); origins != "" {
		cfg.AllowedOrigins = splitAndTrim(origins)
	}

	cfg.DBDSN = buildDatabaseDSN()
	cfg.S3Endpoint = getEnv("BB_S3_ENDPOINT", defaultS3Endpoint)
	cfg.S3Region = getEnv("BB_S3_REGION", defaultS3Region)
	cfg.S3AccessKey = getEnv("BB_S3_ACCESS_KEY", defaultS3AccessKey)
	cfg.S3SecretKey = getEnv("BB_S3_SECRET_KEY", defaultS3SecretKey)
	cfg.S3UseSSL = getBoolEnv("BB_S3_USE_SSL", defaultS3UseSSL)

	validateSecurity(&cfg)

	return cfg
}

func buildDatabaseDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("BB_DB_DSN")); dsn != "" {
		return dsn
	}

	host := getEnv("BB_DB_HOST", defaultDBHost)
	port := getEnv("BB_DB_PORT", defaultDBPort)
	name := getEnv("BB_DB_NAME", defaultDBName)
	user := url.QueryEscape(getEnv("BB_DB_USER", defaultDBUser))
	password := url.QueryEscape(getEnv("BB_DB_PASSWORD", defaultDBPassword))

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return fallback
	}
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	var cleaned []string
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	if len(cleaned) == 0 {
		return []string{"*"}
	}
	return cleaned
}

func containsWildcardOrigin(origins []string) bool {
	for _, origin := range origins {
		if origin == "*" {
			return true
		}
	}
	return false
}

// HasWildcardOrigin reports whether the configuration allows any origin.
func (cfg Config) HasWildcardOrigin() bool {
	return containsWildcardOrigin(cfg.AllowedOrigins)
}

func validateSecurity(cfg *Config) {
	if len(cfg.JWTSecret) < 32 {
		panic("BB_JWT_SECRET must be set to a string with at least 32 characters")
	}
	if cfg.JWTSecret == defaultJWTSecret {
		panic("BB_JWT_SECRET defaults to an insecure value; please override it in the environment")
	}
	if len(cfg.EncryptionKey) != 32 {
		panic("BB_ENCRYPTION_KEY must be exactly 32 bytes for AES-256")
	}
	if string(cfg.EncryptionKey) == defaultEncryptionKey {
		panic("BB_ENCRYPTION_KEY defaults to an insecure value; please override it in the environment")
	}
	if containsWildcardOrigin(cfg.AllowedOrigins) && strings.EqualFold(cfg.Env, "production") {
		panic("BB_ALLOWED_ORIGINS cannot contain '*' when BB_ENV=production")
	}
}
