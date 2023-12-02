package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
	EnvironmentLocal       = "local"
)

type Envs struct {
	Database struct {
		Path string `env:"PATH" env-required:"true"`
		//Host     string `env:"HOST" env-required:"true"`
		//Port     int    `env:"PORT" env-required:"true"`
		//User     string `env:"USER" env-required:"true"`
		//Password string `env:"PASSWORD" env-required:"true"`
		//Database string `env:"DB" env-required:"true"`
	} `env-prefix:"DATABASE_"`
	HTTPServer struct {
		Host        string        `env:"HOST" env-required:"true"`
		Port        int           `env:"PORT" env-required:"true"`
		Address     string        `env:"ADDRESS" env-required:"true"`
		Timeout     time.Duration `env:"TIMEOUT" env-required:"true"`
		IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-required:"true"`
		User        string        `env:"USER" env-required:"true"`
		Password    string        `env:"PASSWORD" env-required:"true"`
	} `env-prefix:"SERVER_"`
	AuthClient struct {
		Address      string        `env:"ADDRESS" env-required:"true"`
		Timeout      time.Duration `env:"TIMEOUT" env-required:"true"`
		RetriesCount int           `env:"RETRIES_COUNT" env-required:"true"`
		Insecure     bool          `env:"INSECURE" env-required:"true"`
	} `env-prefix:"AUTH_CLIENT_"`
	AppSecret string `env:"APP_SECRET" env-required:"true"`
}

var E *Envs
var onceE sync.Once

func ReadEnv(env string) {
	onceE.Do(func() {
		envPath := ".env." + env
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("can't loading env variables, err: %s\n", err.Error())
		}

		log.Printf("[Config] Read %s environment variables\n", env)
		E = &Envs{}
		if err := cleanenv.ReadEnv(E); err != nil {
			help, _ := cleanenv.GetDescription(E, nil)
			log.Println(help)
			log.Fatalln(err)
		}
	})
}

type Config struct {
	Env string `yaml:"environment" env-required:"true"` // local / development / production
}

var C *Config
var onceC sync.Once

func ReadConfig(env string) {
	onceC.Do(func() {
		configPath := fmt.Sprintf("config/config.%s.yaml", env)

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exists: %s, err: %s\n", configPath, err.Error())
		}

		log.Printf("[Config] Read %s configuration variables\n", env)
		C = &Config{}
		if err := cleanenv.ReadConfig(configPath, C); err != nil {
			help, _ := cleanenv.GetDescription(E, nil)
			log.Println(help)
			log.Fatalln(err)
		}
	})
}

func MustLoadVariables() {
	var fenv string
	flag.StringVar(&fenv, "env", "production/ development / local", "project environment")
	flag.Parse()

	ReadEnv(fenv)
	ReadConfig(fenv)
}
