package config

var c *Config

type Config struct {
	StorageDirectory string
	URL              string
}

func Add(config Config) {
	c = &config
}

func Get() *Config {
	return c
}
