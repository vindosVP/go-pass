// Package migratorConfig configures the migrator
package migratorConfig

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// MigratorConfig consists of fields for migrator configuration
type MigratorConfig struct {
	DB             DBConfig `yaml:"db"`
	MigrationsPath string   `yaml:"migrationsPath" validate:"required"`
}

// DBConfig consists of fields for database configuration
type DBConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

// MustLoad loads the MigratorConfig from file
func MustLoad() *MigratorConfig {
	path := configPath()
	if path == "" {
		panic("config path is empty")
	}

	viper.SetConfigFile(path)
	conf := &MigratorConfig{}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	fmt.Println(viper.Get("migrations_path"))

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
