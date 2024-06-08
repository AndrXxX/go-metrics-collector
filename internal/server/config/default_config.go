package config

type Config struct {
	Host string
}

func NewConfig() *Config {
	return &Config{
		Host: "localhost:8080",
	}
}
