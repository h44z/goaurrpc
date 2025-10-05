package config

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"
	"unicode"

	"github.com/kelseyhightower/envconfig"
)

// Duration is a custom duration type that uses seconds instead of nanoseconds as the default unit.
type Duration struct {
	time.Duration
}

func NewDuration(d time.Duration) Duration {
	return Duration{Duration: d}
}

func NewDurationInSeconds(s int) Duration {
	return Duration{Duration: time.Duration(s) * time.Second}
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(d.Seconds()))
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	durationString := strings.Trim(string(b), `"`)
	if durationString == "" {
		return errors.New("time: invalid duration, empty string")
	}

	if unicode.IsDigit(rune(durationString[len(durationString)-1])) {
		durationString += "s" // add default unit if not specified
	}

	d.Duration, err = time.ParseDuration(durationString)
	return
}

// Decode implements envconfig.Decoder
func (d *Duration) Decode(value string) error {
	if value == "" {
		d.Duration = 0
	}

	if unicode.IsDigit(rune(value[len(value)-1])) {
		value += "s" // default unit if none specified
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	d.Duration = parsed
	return nil
}

// Settings is a data structure holding our configuration data
type Settings struct {
	Port                     int      `envconfig:"PORT" default:"10666"`
	AurFileLocation          string   `envconfig:"AUR_FILE_LOCATION" default:"https://aur.archlinux.org/packages-meta-ext-v1.json.gz"`
	MaxResults               int      `envconfig:"MAX_RESULTS" default:"5000"`
	RefreshInterval          Duration `envconfig:"REFRESH_INTERVAL" default:"5m"`
	RateLimit                int      `envconfig:"RATE_LIMIT" default:"4000"`
	LoadFromFile             bool     `envconfig:"LOAD_FROM_FILE" default:"false"`
	RateLimitCleanupInterval Duration `envconfig:"RATE_LIMIT_CLEANUP_INTERVAL" default:"10m"`
	RateLimitTimeWindow      Duration `envconfig:"RATE_LIMIT_TIME_WINDOW" default:"24h"`
	TrustedReverseProxies    []string `envconfig:"TRUSTED_REVERSE_PROXIES" default:"127.0.0.1,::1"`
	EnableSSL                bool     `envconfig:"ENABLE_SSL" default:"false"`
	CertFile                 string   `envconfig:"CERT_FILE" default:""`
	KeyFile                  string   `envconfig:"KEY_FILE" default:""`
	EnableSearchCache        bool     `envconfig:"ENABLE_SEARCH_CACHE" default:"true"`
	CacheCleanupInterval     Duration `envconfig:"CACHE_CLEANUP_INTERVAL" default:"60s"`
	CacheExpirationTime      Duration `envconfig:"CACHE_EXPIRATION_TIME" default:"180s"`
	LogFile                  string   `envconfig:"LOG_FILE" default:""`
	EnableMetrics            bool     `envconfig:"ENABLE_METRICS" default:"true"`
	EnableAdminApi           bool     `envconfig:"ENABLE_ADMIN_API" default:"false"`
	AdminAPIKey              string   `envconfig:"ADMIN_API_KEY" default:"change-me"`
	LogLevel                 string   `envconfig:"LOG_LEVEL" default:"info"` // valid values: debug, info, warn, error
}

// LoadFromEnv processes envconfig annotations and loads environment variables from the environment.
func LoadFromEnv(cfg any) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}

	return nil
}

// Validate config settings
func (s Settings) Validate() error {
	if s.RateLimit == 0 {
		slog.Warn("Rate limiting is disabled - RateLimit = 0")
	}

	errZero := " needs to be specified / greater than 0"
	switch 0 {
	case s.Port:
		return errors.New("config: Port" + errZero)
	case s.MaxResults:
		return errors.New("config: MaxResults" + errZero)
	case int(s.RefreshInterval.Nanoseconds()):
		return errors.New("config: RefreshInterval" + errZero)
	case int(s.RateLimitCleanupInterval.Nanoseconds()):
		return errors.New("config: RateLimitCleanupInterval" + errZero)
	case int(s.RateLimitTimeWindow.Nanoseconds()):
		return errors.New("config: RateLimitTimeWindow" + errZero)
	case int(s.CacheCleanupInterval.Nanoseconds()):
		return errors.New("config: CacheCleanupInterval" + errZero)
	case int(s.CacheExpirationTime.Nanoseconds()):
		return errors.New("config: CacheExpirationTime" + errZero)
	}

	return nil
}

// GetLogLevel returns the correct log level for the set log level string.
func (s Settings) GetLogLevel() slog.Level {
	switch strings.ToLower(s.LogLevel) {
	case "debug", "trace":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warning", "warn":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
