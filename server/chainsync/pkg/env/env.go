package env

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}

func InitEnv() {
	// Determine the build variation (e.g., "staging" or "production")
	buildVariation := "staging"

	// Set up Viper
	viper.SetConfigName(fmt.Sprintf("config.%s", buildVariation))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(getProjectRoot())

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetAlchemyAPIURL() string {
	url := viper.GetString("alchemy.url")
	apiKey := viper.GetString("alchemy.api_key")

	return fmt.Sprintf("%s/%s",
		url, apiKey)
}

func GetServerAddress() string {
	return viper.GetString("server.address")
}

func GetServerPort() string {
	return viper.GetString("server.port")
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
