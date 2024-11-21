package config

import (
	"log"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	Charset  string
}

var DBConfig DatabaseConfig

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	DBConfig = DatabaseConfig{
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		DBName:   viper.GetString("database.dbname"),
		Charset:  viper.GetString("database.charset"),
	}
}
