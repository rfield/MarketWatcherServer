package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	default_grpc_port = 50051
	default_host      = "localhost"
	default_port      = 5432
	default_user      = "postgres"
	default_password  = "admin"
	default_dbname    = "postgres"
)

type Config struct {
	Grpc struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
}

var config Config

func init() {
	_, err := LoadConfig()
	if err != nil {
		panic(err)
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	log.Printf("grpcPort: %d", config.Grpc.Port)

	grpcPort := config.Grpc.Port
	if grpcPort == 0 {
		grpcPort = default_grpc_port
		if os.Getenv("PORT") != "" {
			log.Printf("$PORT = %s", os.Getenv("PORT"))
			var err error
			_, err = fmt.Sscanf(os.Getenv("PORT"), "%d", &grpcPort)
			if err != nil || grpcPort == 0 {
				log.Printf("loadConfig() - invalid PORT environment variable: %v", err)
				grpcPort = default_grpc_port
			}
		}
	}
	config.Grpc.Port = grpcPort

	if config.Database.Host == "" {
		config.Database.Host = default_host
	}
	if config.Database.Port == 0 {
		config.Database.Port = default_port
	}
	if config.Database.User == "" {
		config.Database.User = default_user
	}
	if config.Database.Password == "" {
		config.Database.Password = default_password
	}
	if config.Database.DBName == "" {
		config.Database.DBName = default_dbname
	}

	log.Printf("loadConfig() - loaded config: %+v", config)

	return &config, nil
}

func GetConfig() *Config {
	return &config
}
