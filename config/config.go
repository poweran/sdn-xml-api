package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config структура для представления конфигурационных параметров
type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

// LoadConfig функция для загрузки конфигурационных параметров
func LoadConfig() *Config {
	if _, err := os.Stat("./config/config-local.yml"); err == nil {
		viper.SetConfigName("config-local")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath("./config") // путь к файлу конфигурации

	err := viper.ReadInConfig() // чтение конфигурационного файла
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	viper.AutomaticEnv()      // автоматическое чтение переменных окружения
	viper.SetEnvPrefix("SDN") // префикс для переменных окружения
	viper.BindEnv("DB_HOST")  // связывание переменных окружения с конфигурационными параметрами
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASS")
	viper.BindEnv("DB_NAME")

	config := &Config{
		DBHost: viper.GetString("DB_HOST"),
		DBPort: viper.GetString("DB_PORT"),
		DBUser: viper.GetString("DB_USER"),
		DBPass: viper.GetString("DB_PASS"),
		DBName: viper.GetString("DB_NAME"),
	}

	// Проверяем, что все необходимые параметры указаны в конфигурационном файле
	if config.DBHost == "" || config.DBPort == "" || config.DBUser == "" || config.DBPass == "" || config.DBName == "" {
		log.Fatalf("Missing required config parameters")
	}

	return config
}

// DBConnectionString функция для получения строки подключения к БД
func (c *Config) DBConnectionString() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser + " password=" + c.DBPass + " dbname=" + c.DBName + " sslmode=disable"
}

// GetEnvOrFallback функция для получения значения переменной окружения или её значения по умолчанию
func GetEnvOrFallback(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
