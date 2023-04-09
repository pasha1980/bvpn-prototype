package config

import "bvpn-prototype/internal/protocols/entity"

var c *Config

type Config struct {
	StorageDirectory string
	URL              string
	HttpPort         uint64
	Peers            []entity.Node
}

func Set(config Config) {
	c = &config
}

func Get() *Config {
	return c
}
