package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/canonflow/backend-starter/pkg/helpers"
	"github.com/joho/godotenv"
)

type (
	ConfigKey   int
	ConfigValue interface {
		~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint64 |
			~[]int | ~[]int8 | ~[]int32 | ~[]int64 |
			~string | ~[]string |
			~bool |
			~float32 | ~float64
	}
)

const (
	AppEnv ConfigKey = iota
	AppPort
	AppKey

	BcryptCost

	DBDriver
	DBHost
	DBName
	DBPort
	DBUsername
	DBPassword

	DBIdle
	DBMax
	DBLifeTime

	LogLevel

	RedisHost
	RedisPort
	RedisUsername
	RedisPassword

	CORSAllowOrigins
	CORSAllowCredentials
	CORSAllowHeaders
	CORSAllowMethods

	JWTDurationInMinute
	JWTPath
	JWTDomain
	JWTSecure
	JWTHTTPOnly
)

var cfg *config

type config struct {
	AppEnv  string
	AppPort string
	AppKey  string

	BcryptCost int

	DBDriver   string
	DBHost     string
	DBName     string
	DBPort     string
	DBUsername string
	DBPassword string

	DBIdle     int
	DBMax      int
	DBLifetime int

	LogLevel string

	RedisHost     string
	RedisPort     string
	RedisUsername string
	RedisPassword string

	CORSAllowOrigins     []string
	CORSAllowCredentials bool
	CORSAllowHeaders     []string
	CORSAllowMethods     []string

	JWTDurationInMinute uint
	JWTPath             string
	JWTDomain           string
	JWTSecure           bool
	JWTHTTPOnly         bool
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load .env file: %v", err.Error()))
	}

	loadConfig()
}

func loadConfig() {
	if cfg == nil {
		cfg = &config{}

		// App Config
		cfg.AppEnv = os.Getenv("APP_ENV")
		cfg.AppPort = os.Getenv("APP_PORT")
		cfg.AppKey = os.Getenv("APP_KEY")
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
		cfg.BcryptCost = helpers.Parser(
			strconv.Atoi,
			os.Getenv("BCRYPT_COST"),
			10,
		)

		// Database
		cfg.DBDriver = os.Getenv("DB_DRIVER")
		cfg.DBHost = os.Getenv("DB_HOST")
		cfg.DBPort = os.Getenv("DB_PORT")
		cfg.DBName = os.Getenv("DB_NAME")
		cfg.DBUsername = os.Getenv("DB_USERNAME")
		cfg.DBPassword = os.Getenv("DB_PASSWORD")
		cfg.DBIdle = helpers.Parser(
			strconv.Atoi,
			os.Getenv("DB_IDLE"),
			10,
		)
		cfg.DBMax = helpers.Parser(
			strconv.Atoi,
			os.Getenv("DB_MAX"),
			100,
		)
		cfg.DBLifetime = helpers.Parser(
			strconv.Atoi,
			os.Getenv("DB_LIFETIME"),
			300,
		)

		// Redis
		cfg.RedisHost = os.Getenv("REDIS_HOST")
		cfg.RedisPort = os.Getenv("REDIS_PORT")
		cfg.RedisUsername = os.Getenv("REDIS_USERNAME")
		cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")

		// CORS
		cfg.CORSAllowOrigins = helpers.SplitStringTrim(
			os.Getenv("CORS_ALLOW_ORIGINS"),
			",",
		)
		cfg.CORSAllowCredentials = helpers.Parser(
			strconv.ParseBool,
			os.Getenv("CORS_ALLOW_CREDENTIALS"),
			true,
		)
		cfg.CORSAllowHeaders = helpers.SplitStringTrim(
			os.Getenv("CORS_ALLOW_HEADERS"),
			",",
		)
		cfg.CORSAllowMethods = helpers.SplitStringTrim(
			os.Getenv("CORS_ALLOW_METHODS"),
			",",
		)

		// JWT
		cfg.JWTDurationInMinute = helpers.Parser(
			helpers.ParseUint,
			os.Getenv("JWT_DURATION_IN_MINUTES"),
			3600,
		)
		cfg.JWTPath = os.Getenv("JWT_PATH")
		cfg.JWTDomain = os.Getenv("JWT_DOMAIN")
		cfg.JWTSecure = helpers.Parser(
			strconv.ParseBool,
			os.Getenv("JWT_SECURE"),
			false,
		)
		cfg.JWTHTTPOnly = helpers.Parser(
			strconv.ParseBool,
			os.Getenv("JWT_HTTP_ONLY"),
			true,
		)
	}
}

func (c *config) get(key ConfigKey) any {
	switch key {
	case AppEnv:
		return c.AppEnv
	case AppPort:
		return c.AppPort
	case AppKey:
		return c.AppKey
	case BcryptCost:
		return c.BcryptCost
	case DBDriver:
		return c.DBDriver
	case DBHost:
		return c.DBHost
	case DBPort:
		return c.DBPort
	case DBName:
		return c.DBName
	case DBUsername:
		return c.DBUsername
	case DBPassword:
		return c.DBPassword
	case DBIdle:
		return c.DBIdle
	case DBMax:
		return c.DBMax
	case DBLifeTime:
		return c.DBLifetime
	case LogLevel:
		return c.LogLevel
	case RedisHost:
		return c.RedisHost
	case RedisPort:
		return c.RedisPort
	case RedisUsername:
		return c.RedisUsername
	case RedisPassword:
		return c.RedisPassword
	case CORSAllowOrigins:
		return c.CORSAllowOrigins
	case CORSAllowCredentials:
		return c.CORSAllowCredentials
	case CORSAllowHeaders:
		return c.CORSAllowHeaders
	case CORSAllowMethods:
		return c.CORSAllowMethods
	case JWTDurationInMinute:
		return c.JWTDurationInMinute
	case JWTPath:
		return c.JWTPath
	case JWTDomain:
		return c.JWTDomain
	case JWTSecure:
		return c.JWTSecure
	case JWTHTTPOnly:
		return c.JWTHTTPOnly
	default:
		return nil
	}
}

// Get returns the value or the zero value of T if missing/wrong type.
func Get[T ConfigValue](key ConfigKey) T {
	loadConfig()

	v := cfg.get(key)
	t, _ := v.(T)
	return t
}

// GetOrDefault returns the value, or def if missing/wrong type.
func GetOrDefault[T ConfigValue](key ConfigKey, def T) T {
	loadConfig()

	v := cfg.get(key)
	t, ok := v.(T)
	if !ok {
		return def
	}
	return t
}
