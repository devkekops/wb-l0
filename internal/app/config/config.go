package config

type Config struct {
	RunAddress  string `env:"RUN_ADDRESS"`
	DatabaseURI string `env:"DATABASE_URI"`
	NatsURI     string `env:"NATS_URI"`
}
