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
	ApiToken           string
	AppHost            string
	DbConnectionString string
}

type RawColumnsConfig struct {
	Columns map[string]struct {
		Description string `yaml:"description"`
		Unit        string `yaml:"unit"`
	} `yaml:"columns"`
}

type ColumnsConfig struct {
	DateFormat        string
	HumidityFormat    string
	TemperatureFormat string
}

var (
	conf          *Config
	onceConfigs   sync.Once
	onceColumns   sync.Once
	columnsConfig *ColumnsConfig
)

var dateFormats = map[string]string{
	"YYYY-MM-DD": "2006-01-02",
}

func GetColumns() *ColumnsConfig {
	onceColumns.Do(func() {
		columnsYaml, err := os.ReadFile("configs/columns.yaml")
		if err != nil {
			log.Fatalln(err)
		}

		var rawColumnsConfig = RawColumnsConfig{}

		if err := yaml.Unmarshal(columnsYaml, &rawColumnsConfig); err != nil {
			log.Fatalln(err)
		}
		dateFormat := dateFormats[rawColumnsConfig.Columns["Date"].Unit]
		if dateFormat == "" {
			log.Fatalln("invalid date format specified in columns.yaml")
		}

		columnsConfig = &ColumnsConfig{
			DateFormat:        dateFormat,
			HumidityFormat:    rawColumnsConfig.Columns["Humidity"].Unit,
			TemperatureFormat: rawColumnsConfig.Columns["Temperature"].Unit,
		}
	})
	return columnsConfig
}

func Get() *Config {
	onceConfigs.Do(func() {
		appEnv := os.Getenv("APP_ENV")
		if appEnv == "" {
			log.Fatalln("APP_ENV environment variable is not set")
		}

		envFilename := fmt.Sprintf("./configs/.env.%s", appEnv)
		err := godotenv.Load(envFilename)
		if err != nil {
			log.Printf("No %s file found. Using system environment variables\n", envFilename)
		}

		conf = &Config{
			ApiToken:           os.Getenv("API_TOKEN"),
			AppHost:            os.Getenv("APP_HOST"),
			DbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		}

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
	Get()
	GetColumns()
}
