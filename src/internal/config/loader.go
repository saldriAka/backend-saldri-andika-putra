package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Get() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error when loading file config: ", err.Error())
	}

	expint, _ := strconv.Atoi(os.Getenv("JWT_EXP"))

	return &Config{
		Server: Server{
			Host:   os.Getenv("SERVER_HOST"),
			Port:   os.Getenv("SERVER_PORT"),
			Assets: os.Getenv("SERVER_ASSETS_URL"),
		},
		Database: Database{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
			Tz:   os.Getenv("DB_TZ"),
		},
		Jwt: Jwt{
			Key: os.Getenv("JWT_KEY"),
			Exp: expint,
		},
		Storage: Storage{
			BasePath: os.Getenv("STORAGE_PATH"),
		},
		ServiceMode: ServiceMode{
			ServiceMode: os.Getenv("SERVICE_MODE"),
		},
	}
}
