package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database Database
	Kafka    Kafka
	Cache    Cache
	Port     string `envconfig:"PORT" default:":8080"`
}

type Database struct {
	Host              string `envconfig:"DB_HOST" default:"localhost"`
	Port              string `envconfig:"DB_PORT" required:"true"`
	User              string `envconfig:"DB_USER" required:"true"`
	Password          string `envconfig:"DB_PASSWORD" required:"true"`
	Name              string `envconfig:"DB_NAME" required:"true"`
	Schema            string `envconfig:"DB_SCHEMA" required:"true" default:"orders"`
	MaxOpenConnection int    `envconfig:"DB_MAX_OPEN_CONNECTION" default:"10"`
}

type Kafka struct {
	Host    string `envconfig:"KAFKA_HOST" required:"true"`
	Port    string `envconfig:"KAFKA_PORT" required:"true"`
	Retry   int    `envconfig:"KAFKA_RETRY"  default:"2"`
	Backoff int    `envconfig:"KAFKA_BACKOFF"  default:"100"`
}

type Cache struct {
	Size int `envconfig:"CACHE_SIZE" required:"true"`
	TTL  int `envconfig:"CACHE_TTL" required:"true"`
}

func NewParsedConfig() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
