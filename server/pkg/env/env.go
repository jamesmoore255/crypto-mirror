package env

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitEnv() {
	// Determine the build variation (e.g., "staging" or "production")
	buildVariation := "staging"

	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigFile(fmt.Sprintf("config.%s.yaml", buildVariation))

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetServerAddress() string {
	return viper.GetString("server.address")
}

func GetServerPort() int {
	return viper.GetInt("server.port")
}

func GetAlchemyAPIKey() string {
	return viper.GetString("alchemy.api_key")
}

func GetDatabaseURL() string {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.database_name")

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbUsername, dbPassword, dbHost, dbPort, dbName)
}
