package config

import (
	"fmt"
	"log"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)


const (
	cfgFileName = "dev"
	cfgFileType = "env"
)


type Config struct {
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"  validate:"required"`

	DBSource string `mapstructure:"DB_SOURCE" validate:"required"` 

	PgHost     string `mapstructure:"PG_HOST" validate:"required"`
	PgPort     string `mapstructure:"PG_PORT" validate:"required"`
	PgUser     string `mapstructure:"PG_USER" validate:"required"`
	PgPassword string `mapstructure:"PG_PASSWORD" validate:"required"`
	PgDB       string `mapstructure:"PG_DB" validate:"required"`
	PgSSL      string `mapstructure:"PG_SSL" validate:"required"`
}

// Load reads configuration from file or environment variables.
func Load() (config Config, err error) {
	return LoadFromPath(".")
}


func LoadFromPath(cfgPath string) (config Config, err error) {
	viper.AddConfigPath(cfgPath)
	viper.SetConfigName(cfgFileName)
	viper.SetConfigType(cfgFileType)

	// bind env vars

	_ = viper.BindEnv("SERVER_ADDRESS")

	_ = viper.BindEnv("DB_SOURCE")

	_ = viper.BindEnv("PG_HOST")
	_ = viper.BindEnv("PG_PORT")
	_ = viper.BindEnv("PG_USER")
	_ = viper.BindEnv("PG_PASSWORD")
	_ = viper.BindEnv("PG_DB")
	_ = viper.BindEnv("PG_SSL")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("config file '%s' not found", path.Join(cfgPath, fmt.Sprintf("%s.%s", cfgFileName, cfgFileType)))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	v := validator.New()
	err = v.Struct(&config)
	return config, err
}