package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RpcURL     string
	ServerPort string
	WebPort    string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env")
	}

	return &Config{
		RpcURL:     os.Getenv("RPC_URL"),
		ServerPort: os.Getenv("SERVER_PORT"),
		WebPort:    os.Getenv("WEB_PORT"),
	}
}
