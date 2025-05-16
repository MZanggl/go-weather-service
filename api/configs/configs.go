package configs

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	ApiToken           string
	AppHost            string
	DbConnectionString string
}

var (
	conf *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found. Using system environment variables")
		}

		conf = &Config{
			ApiToken:           os.Getenv("API_TOKEN"),
			AppHost:            os.Getenv("APP_HOST"),
			DbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		}
		log.Println("Loading configuration...", conf)

		if conf.ApiToken == "" {
			log.Fatalln("API_TOKEN environment variable is not set")
		}
		if conf.AppHost == "" {
			log.Fatalln("APP_HOST environment variable is not set")
		}
		if conf.DbConnectionString == "" {
			log.Fatalln("DB_CONNECTION_STRING environment variable is not set")
		}
		log.Println("Configuration loaded successfully")
	})
	return conf
}
