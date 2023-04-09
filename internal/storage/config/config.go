package config

import "bvpn-prototype/internal/protocols/entity"

var c *Config

type Config struct {
	StorageDirectory string
	URL              string
	HttpPort         string
	VpnPort          string
	VpnProto         string
	Peers            []entity.Node
}

func Set(config Config) {
	c = &config
}

func Get() *Config {
	return c
}
