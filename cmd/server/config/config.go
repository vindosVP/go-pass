// Package serverConfig configures the server
package serverConfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

// ServerConfig consists of fields for server configuration
type ServerConfig struct {
	Env  string     `yaml:"env" validate:"required"`
	DB   DBConfig   `yaml:"db"`
	GRPC GRPCConfig `yaml:"grpc"`
	Auth AuthConfig `yaml:"auth"`
}

// String turns the ServerConfig to string
func (s *ServerConfig) String() string {
	out, err := json.Marshal(s)

	if err != nil {
		sl.Err(fmt.Errorf("failed to marshal server config: %w", err))
	}

	return string(out)
}

// DBConfig consists of fields for database configuration
type DBConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

// GRPCConfig consists of fields for grpc configuration
type GRPCConfig struct {
	Port    int           `yaml:"port" validate:"required"`
	Timeout time.Duration `yaml:"timeout" validate:"required"`
}

// AuthConfig consists of fields for authentication configuration
type AuthConfig struct {
	TokenTTL time.Duration `yaml:"tokenTTL" validate:"required"`
}

// MustLoad loads the ServerConfig from file
func MustLoad() *ServerConfig {
	path := configPath()
	if path == "" {
		panic("config path is empty")
	}

	viper.SetConfigFile(path)
	conf := &ServerConfig{}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
		panic(fmt.Errorf("missing requiered attributes: %w", err))
	}

	return conf
}

func configPath() string {
	var res string

	flag.StringVar(&res, "c", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
