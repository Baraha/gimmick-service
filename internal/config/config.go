package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	"log"
	"os"
	"time"
)

const (
	HeaderContentTypeKey  = "Content-Type"
	HeaderContentTypeJSON = "application/json; charset=utf-8"
)

func NewConfig() (*Configuration, error) {
	var envFiles []string
	if _, err := os.Stat(".env"); err == nil {
		log.Println("found .env file, adding it to env config files list")
		envFiles = append(envFiles, ".env")
	}
	if os.Getenv("APP_ENV") != "" {
		appEnvName := fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))
		if _, err := os.Stat(appEnvName); err == nil {
			log.Println("found", appEnvName, "file, adding it to env config files list")
			envFiles = append(envFiles, appEnvName)
		}
	}
	if len(envFiles) > 0 {
		err := godotenv.Overload(envFiles...)
		if err != nil {
			return nil, errors.Wrapf(err, "error while opening env config: %s", err)
		}
	}
	cfg := &Configuration{}
	ctx := context.Background()

	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "error while parsing env config: %s", err)
	}
	return cfg, nil
}

// Configuration is basic structure that contains configuration
type Configuration struct {
	HTTPServer  HTTPServerConfig `env:",prefix=HTTP_"`
	FileStorage FS               `env:",prefix=FILE_STORAGE_"`
	GRPCServer  GRPCServerConfig `env:",prefix=GRPC_SERVER_"`
}

type GRPCServerConfig struct {
	RequestLoggingEnabled      bool          `env:"REQUEST_LOGGING_ENABLED,default=true"`
	ResponseTimeLoggingEnabled bool          `env:"RESPONSE_TIME_LOGGING_ENABLED,default=false"`
	ReadTimeout                time.Duration `env:"READ_TIMEOUT,default=30s"`
	WriteTimeout               time.Duration `env:"WRITE_TIMEOUT,default=30s"`
	IdleTimeout                time.Duration `env:"IDLE_TIMEOUT,default=30s"`
	MaxRequestBodySize         int           `env:"MAX_REQUEST_BODY_SIZE,default=4194304"`
	Network                    string        `env:"NETWORK,default=tcp"`
	Address                    string        `env:"ADDRESS,default=:18089"`
}

type HTTPServerConfig struct {
	CORSEnabled                bool          `env:"CORS_ENABLED,default=false"`
	RequestLoggingEnabled      bool          `env:"REQUEST_LOGGING_ENABLED,default=true"`
	ResponseTimeLoggingEnabled bool          `env:"RESPONSE_TIME_LOGGING_ENABLED,default=false"`
	ReadTimeout                time.Duration `env:"READ_TIMEOUT,default=30s"`
	WriteTimeout               time.Duration `env:"WRITE_TIMEOUT,default=30s"`
	IdleTimeout                time.Duration `env:"IDLE_TIMEOUT,default=30s"`
	MaxRequestBodySize         int           `env:"MAX_REQUEST_BODY_SIZE,default=4194304"`
	Network                    string        `env:"NETWORK,default=tcp"`
	Address                    string        `env:"ADDRESS,default=:9090"`
}

type FS struct {
	DefaultDir string `env:"LEVEL,default=./"`
}
