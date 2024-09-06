package configs

import (
	"fmt"
	"os"

	"github.com/spf13/cast"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	MongoHost     string
	MongoPort     string
	MongoUser     string
	MongoPassword string
	MongoDB       string

	ServiceName string
	Environment string
	LoggerLevel string

	TaskServiceGrpcHost string
	TaskServiceGrpcPort string
	EmailPassword       string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found", err)
	}

	cfg := Config{}

	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	cfg.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", 5433))
	cfg.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1111"))
	cfg.PostgresDB = cast.ToString(getOrReturnDefault("POSTGRES_DB", "task_service"))

	cfg.MongoHost = cast.ToString(getOrReturnDefault("MONGO_HOST", "localhost"))
	cfg.MongoPort = cast.ToString(getOrReturnDefault("MONGO_PORT", 27017))
	cfg.MongoUser = cast.ToString(getOrReturnDefault("MONGO_USER", "mongosh"))
	cfg.MongoPassword = cast.ToString(getOrReturnDefault("MONGO_PASSWORD", "1111"))
	cfg.MongoDB = cast.ToString(getOrReturnDefault("MONGO_DB", "task_service"))

	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "task_service"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))

	cfg.TaskServiceGrpcHost = cast.ToString(getOrReturnDefault("TASK_SERVICE_GRPC_HOST", "localhost"))
	cfg.TaskServiceGrpcPort = cast.ToString(getOrReturnDefault("TASK_SERVICE_GRPC_PORT", ":8082"))

	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
