package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
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
	// 加载 .env 文件
	err := godotenv.Load("./configs/.env")
	if err != nil {
		// 没读到就忽略，但是log
		log.Printf("Error loading .env file, %s", err)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	DBConfig = DatabaseConfig{
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		DBName:   viper.GetString("database.dbname"),
		Charset:  viper.GetString("database.charset"),
	}
}
