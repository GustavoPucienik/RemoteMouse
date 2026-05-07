package config

import "fmt"

type Config struct {
	Host string
	Port int
}

func Default() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: 8080,
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
