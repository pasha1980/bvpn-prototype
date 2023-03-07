package main

import (
	"bvpn-prototype/internal"
	"bvpn-prototype/internal/protocols/entity"
	"fmt"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const defaultConfigFile = "/etc/bvpn/config.yaml"

var configFile = defaultConfigFile

func main() {
	var opts struct {
		ConfigFile string `short:"c" long:"configs" description:"Configuration file" required:"false"`
	}

	commands, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err)
	}
	command := commands[0]
	configFile = opts.ConfigFile

	switch command {
	case "run":
		kernel, err := kernelFromConfigs()
		if err != nil {
			log.Fatalln(err)
		}

		kernel.Run()
	default:
		fmt.Println("Hello") // todo
	}
}

func kernelFromConfigs() (*internal.Kernel, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg struct {
		//Price struct {
		//	Gb   *float64 `yaml:"gb"`
		//	Mb   *float64 `yaml:"mb"`
		//	Kb   *float64 `yaml:"kb"`
		//	Bvpn float64  `yaml:"bvpn"`
		//} `yaml:"price"`
		HttpUrl string `yaml:"http_url"`
		Ports   struct {
			//Vpn  uint64 `yaml:"vpn"`
			Http uint64 `yaml:"http"`
		} `yaml:"ports"`
		Peers []struct {
			Ip string `yaml:"ip"`
			//VpnUrl  string  `yaml:"vpn_url"`
			HttpUrl string `yaml:"http_url"`
			//Secret  *string `yaml:"secret"`
		} `yaml:"peers"`
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	var peers []entity.Node
	for _, peerCfg := range cfg.Peers {
		peers = append(peers, entity.Node{
			URL: peerCfg.HttpUrl,
			IP:  peerCfg.Ip,
		})
	}

	if cfg.Ports.Http == 0 {
		cfg.Ports.Http = 80
	}

	kernel := internal.Kernel{
		URL:      cfg.HttpUrl,
		HttpPort: cfg.Ports.Http,
		Peers:    peers,
	}

	return &kernel, nil
}
