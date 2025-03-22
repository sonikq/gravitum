package config

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	RunAddress string

	DatabaseDSN   string
	DBPoolWorkers int

	CtxTimeOut  time.Duration
	LogLevel    string
	ServiceName string
}

const (
	defaultRunAddress    = "localhost:3000"
	defaultLogLevel      = "info"
	defaultServiceName   = "user-management"
	defaultDatabaseDSN   = ""
	defaultDBPoolWorkers = 50
	defaultCtxTimeOut    = 5 * time.Second
)

func Load(envFiles ...string) (Config, error) {
	if len(envFiles) != 0 {
		if err := godotenv.Load(envFiles...); err != nil {
			return Config{}, err
		}
	}

	runAddress := flag.String("a", defaultRunAddress,
		"run address defines on what port and host the server will be started")
	databaseDSN := flag.String("d", defaultDatabaseDSN,
		"defines the database connection address")

	cfg := new(Config)

	cfg.RunAddress = getEnvString(*runAddress, "RUN_ADDRESS")

	cfg.DatabaseDSN = getEnvString(*databaseDSN, "DATABASE_DSN")
	cfg.DBPoolWorkers = getEnvInt(defaultDBPoolWorkers, "DB_POOL_WORKERS")

	cfg.CtxTimeOut = getEnvDuration(defaultCtxTimeOut, "CTX_TIMEOUT")
	cfg.LogLevel = getEnvString(defaultLogLevel, "LOG_LEVEL")
	cfg.ServiceName = getEnvString(defaultServiceName, "SERVICE_NAME")

	return *cfg, nil
}

func getEnvString(flagValue string, envKey string) string {
	envValue, exists := os.LookupEnv(envKey)
	if exists {
		return envValue
	}
	return flagValue
}

func getEnvDuration(flagValue time.Duration, envKey string) time.Duration {
	envValue, exists := os.LookupEnv(envKey)
	if exists {
		return time.Millisecond * time.Duration(cast.ToInt(envValue))
	}
	return flagValue
}

func getEnvInt(flagValue int, envKey string) int {
	envValue, exists := os.LookupEnv(envKey)
	if exists {
		intVal, err := strconv.Atoi(envValue)
		if err != nil {
			log.Printf("cant convert env-key: %s to int", envValue)
			return 1
		}

		return intVal
	}

	return flagValue
}
