package common

import (
	"github.com/spf13/viper"
)

// Configuration keep all configuration
type Configuration struct {
	DataBaseConfig *DataBaseConfig
	ServerConfig   *ServerConfig
}

// DataBaseConfig represents the Database configurations.
type DataBaseConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	Timeout        uint
	MinConnections int
	MaxConnections int
}

type BrokerConfig struct {
	URLs []string
}

// ServerConfig represents server configuration to initialize.
type ServerConfig struct {
	Port string
}

// CONFIGURATION initializes a new Configuration
func NewConfiguration() *Configuration {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// set default values
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_TIMEOUT", 10)
	viper.SetDefault("DB_MIN_CONNECTIONS", 1)
	viper.SetDefault("DB_MAX_CONNECTIONS", 5)

	// environment variables
	return &Configuration{
		DataBaseConfig: &DataBaseConfig{
			Host:           viper.GetString("DB_HOST"),
			Port:           viper.GetString("DB_PORT"),
			User:           viper.GetString("DB_USER"),
			Password:       viper.GetString("DB_PASSWORD"),
			Name:           viper.GetString("DB_NAME"),
			Timeout:        viper.GetUint("DB_TIMEOUT"),
			MinConnections: viper.GetInt("DB_MIN_CONNECTIONS"),
			MaxConnections: viper.GetInt("DB_MAX_CONNECTIONS"),
		},
		ServerConfig: &ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
	}
}
