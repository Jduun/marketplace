package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort         int    `env:"APP_PORT"`
	JWTSecret       string `env:"JWT_SECRET"`
	AppEnv          AppEnv `env:"APP_ENV"`
	TokenTTLMinutes int    `env:"TOKEN_TTL_MINUTES"`
	DBHost          string `env:"DB_HOST"`
	DBPort          int    `env:"DB_PORT"`
	DBUsername      string `env:"DB_USERNAME"`
	DBPassword      string `env:"DB_PASSWORD"`
	DBName          string `env:"DB_NAME"`
	DBPath          string `env:"DB_PATH"`
}

type AppEnv string

const (
	Local AppEnv = "local"
	Dev   AppEnv = "dev"
	Prod  AppEnv = "prod"
)

var (
	once sync.Once
	Cfg  *Config
)

func MustLoad() *Config {
	once.Do(func() {
		Cfg = &Config{}
		if err := cleanenv.ReadEnv(Cfg); err != nil {
			log.Fatalf("Cannot read .env file: %s", err)
		}
		fmt.Println(fmt.Sprintf("APP_PORT: %s", Cfg.AppPort))
	})
	return Cfg
}

func (cfg *Config) GetDBURL() string {
	DBURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBName,
	)
	return DBURL
}
