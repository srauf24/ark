package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	Integration   IntegrationConfig    `koanf:"integration" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

type IntegrationConfig struct {
	ResendAPIKey string `koanf:"resend_api_key" validate:"required"`
}

type AuthConfig struct {
	SecretKey string      `koanf:"secret_key" validate:"required"`
	Clerk     ClerkConfig `koanf:"clerk" validate:"required"`
}

type ClerkConfig struct {
	SecretKey    string `koanf:"secret_key" validate:"required"`
	JWTIssuer    string `koanf:"jwt_issuer" validate:"required,url"`
	PEMPublicKey string `koanf:"pem_public_key"`
}

func parseMapString(value string) (map[string]string, bool) {
	if !strings.HasPrefix(value, "map[") || !strings.HasSuffix(value, "]") {
		return nil, false
	}

	content := strings.TrimPrefix(value, "map[")
	content = strings.TrimSuffix(content, "]")

	if content == "" {
		return make(map[string]string), true
	}

	result := make(map[string]string)
	i := 0

	for i < len(content) {
		keyStart := i
		for i < len(content) && content[i] != ':' {
			i++
		}
		if i >= len(content) {
			break
		}

		key := strings.TrimSpace(content[keyStart:i])
		i++

		valueStart := i

		// detect nested map
		if i+4 <= len(content) && content[i:i+4] == "map[" {
			bracketCount := 0
			for i < len(content) {
				if i+4 <= len(content) && content[i:i+4] == "map[" {
					bracketCount++
					i += 4
				} else if content[i] == ']' {
					bracketCount--
					i++
					if bracketCount == 0 {
						break
					}
				} else {
					i++
				}
			}
		} else {
			for i < len(content) && content[i] != ' ' {
				i++
			}
		}

		v := strings.TrimSpace(content[valueStart:i])

		if nested, ok := parseMapString(v); ok {
			for mk, mv := range nested {
				result[key+"."+mk] = mv
			}
		} else {
			result[key] = v
		}

		for i < len(content) && content[i] == ' ' {
			i++
		}
	}

	return result, true
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	// Load ARK_* env vars into koanf
	err := k.Load(env.ProviderWithValue("ARK_", ".", func(key, value string) (string, any) {
		// Clean up the key: ARK_REDIS -> redis
		k := strings.ToLower(strings.TrimPrefix(key, "ARK_"))

		// Special handling for CORS origins - split comma-separated values
		if k == "server.cors_allowed_origins" {
			origins := strings.Split(value, ",")
			trimmedOrigins := make([]string, len(origins))
			for i, origin := range origins {
				trimmedOrigins[i] = strings.TrimSpace(origin)
			}
			return k, trimmedOrigins
		}

		// Check if it's a map and parse it directly here
		if mapData, isMap := parseMapString(value); isMap {
			return k, mapData
		}

		// Otherwise return the raw value
		return k, value
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial ARK_ environment variables")
	}

	mainConfig := &Config{}

	if err := k.Unmarshal("", mainConfig); err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal config")
	}

	validate := validator.New()

	if err := validate.Struct(mainConfig); err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	} else {
		logger.Info().Msg("config validation passed")
	}

	// Apply defaults
	defaults := DefaultObservabilityConfig()
	if mainConfig.Observability == nil {
		mainConfig.Observability = defaults
	} else {
		// Merge defaults for missing fields
		if mainConfig.Observability.Logging.Level == "" {
			mainConfig.Observability.Logging.Level = defaults.Logging.Level
		}
		if mainConfig.Observability.Logging.Format == "" {
			mainConfig.Observability.Logging.Format = defaults.Logging.Format
		}
		if mainConfig.Observability.Logging.SlowQueryThreshold == 0 {
			mainConfig.Observability.Logging.SlowQueryThreshold = defaults.Logging.SlowQueryThreshold
		}
		// Ensure HealthChecks defaults are applied if not set
		if mainConfig.Observability.HealthChecks.Interval == 0 {
			mainConfig.Observability.HealthChecks.Interval = defaults.HealthChecks.Interval
		}
		if mainConfig.Observability.HealthChecks.Timeout == 0 {
			mainConfig.Observability.HealthChecks.Timeout = defaults.HealthChecks.Timeout
		}
	}

	// Override service name and environment from primary config
	mainConfig.Observability.ServiceName = "ark"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	// Validate observability config
	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	return mainConfig, nil
}
