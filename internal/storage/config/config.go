package config

var c *Config

type Config struct {
	StorageDirectory string
	URL              string
}

func Set(config Config) {
	c = &config
}

func Get() *Config {
	return c
}
