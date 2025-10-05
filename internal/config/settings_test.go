package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		expected   Settings
		shouldFail bool
	}{
		{
			name:    "Default configuration",
			envVars: map[string]string{},
			expected: Settings{
				Port:                     10666,
				AurFileLocation:          "https://aur.archlinux.org/packages-meta-ext-v1.json.gz",
				MaxResults:               5000,
				RefreshInterval:          NewDuration(5 * time.Minute),
				RateLimit:                4000,
				LoadFromFile:             false,
				RateLimitCleanupInterval: NewDuration(10 * time.Minute),
				RateLimitTimeWindow:      NewDuration(24 * time.Hour),
				TrustedReverseProxies:    []string{"127.0.0.1", "::1"},
				EnableSSL:                false,
				CertFile:                 "",
				KeyFile:                  "",
				EnableSearchCache:        true,
				CacheCleanupInterval:     NewDuration(60 * time.Second),
				CacheExpirationTime:      NewDuration(180 * time.Second),
				LogFile:                  "",
				EnableMetrics:            true,
				EnableAdminApi:           false,
				AdminAPIKey:              "change-me",
				LogLevel:                 "info",
			},
			shouldFail: false,
		},
		{
			name: "Valid configuration overrides",
			envVars: map[string]string{
				"PORT":                        "8080",
				"MAX_RESULTS":                 "100",
				"REFRESH_INTERVAL":            "30s",
				"RATE_LIMIT_CLEANUP_INTERVAL": "1m15s",
				"RATE_LIMIT_TIME_WINDOW":      "24h",
				"CACHE_CLEANUP_INTERVAL":      "20us",
				"CACHE_EXPIRATION_TIME":       "60", // implicitly converted to seconds
			},
			expected: Settings{
				Port:                     8080,
				AurFileLocation:          "https://aur.archlinux.org/packages-meta-ext-v1.json.gz",
				MaxResults:               100,
				RefreshInterval:          NewDuration(30 * time.Second),
				RateLimit:                4000,
				LoadFromFile:             false,
				RateLimitCleanupInterval: NewDuration(1*time.Minute + 15*time.Second),
				RateLimitTimeWindow:      NewDuration(24 * time.Hour),
				TrustedReverseProxies:    []string{"127.0.0.1", "::1"},
				EnableSSL:                false,
				CertFile:                 "",
				KeyFile:                  "",
				EnableSearchCache:        true,
				CacheCleanupInterval:     NewDuration(20 * time.Microsecond),
				CacheExpirationTime:      NewDuration(60 * time.Second),
				LogFile:                  "",
				EnableMetrics:            true,
				EnableAdminApi:           false,
				AdminAPIKey:              "change-me",
				LogLevel:                 "info",
			},
			shouldFail: false,
		},
		{
			name: "Invalid field type",
			envVars: map[string]string{
				"PORT": "invalid",
			},
			expected:   Settings{},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range test.envVars {
				t.Setenv(key, value)
			}

			var s Settings
			err := LoadFromEnv(&s)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, s)
			}
		})
	}
}

func TestValidateSettings(t *testing.T) {
	var s Settings

	// load defaults
	err := LoadFromEnv(&s)
	require.NoError(t, err)

	// ok
	err = s.Validate()
	assert.Nil(t, err)

	// errors
	s = Settings{} // empty, uninitialized settings
	err = s.Validate()
	assert.Error(t, err)

	s.Port = 1
	err = s.Validate()
	assert.Error(t, err)

	s.MaxResults = 1
	err = s.Validate()
	assert.Error(t, err)

	s.RefreshInterval = NewDurationInSeconds(1)
	err = s.Validate()
	assert.Error(t, err)

	s.RateLimitCleanupInterval = NewDurationInSeconds(1)
	err = s.Validate()
	assert.Error(t, err)

	s.RateLimitTimeWindow = NewDurationInSeconds(1)
	err = s.Validate()
	assert.Error(t, err)

	s.CacheCleanupInterval = NewDurationInSeconds(1)
	err = s.Validate()
	assert.Error(t, err)

	s.CacheExpirationTime = NewDurationInSeconds(1)
	err = s.Validate()
	assert.NoError(t, err)
}
