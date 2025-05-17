package configs

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	ApiToken           string
	AppHost            string
	DbConnectionString string
}

type ColumnsConfig struct {
	Columns map[string]struct {
		Description string `yaml:"description"`
		Unit        string `yaml:"unit"`
	} `yaml:"columns"`
}

var (
	conf          *Config
	onceConfigs   sync.Once
	onceColumns   sync.Once
	columnsConfig *ColumnsConfig
)

func GetColumns() *ColumnsConfig {
	onceColumns.Do(func() {
		columnsYaml, err := os.ReadFile("configs/columns.yaml")
		if err != nil {
			log.Fatalln(err)
		}
		columnsConfig = &ColumnsConfig{}
		if err := yaml.Unmarshal(columnsYaml, columnsConfig); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Columns config loaded successfully", columnsConfig)
	})
	return columnsConfig
}

func Load() *Config {
	onceConfigs.Do(func() {
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

func init() {
	// Load the configuration once during startup
	Load()
	GetColumns()
}
