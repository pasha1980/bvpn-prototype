package main

import (
	"bvpn-prototype/internal"
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
		fmt.Println("Hello")
	}
}

func kernelFromConfigs() (*internal.Kernel, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg struct {
		Wallet string `yaml:"wallet"`
		Price  struct {
			Gb   *float64 `yaml:"gb"`
			Mb   *float64 `yaml:"mb"`
			Kb   *float64 `yaml:"kb"`
			Bvpn float64  `yaml:"bvpn"`
		} `yaml:"price"`
		Ports struct {
			Vpn  uint `yaml:"vpn"`
			Http uint `yaml:"http"`
		} `yaml:"ports"`
		Peers []struct {
			Ip      string  `yaml:"ip"`
			VpnUrl  string  `yaml:"vpn_url"`
			HttpUrl string  `yaml:"http_url"`
			Secret  *string `yaml:"secret"`
		} `yaml:"peers"`
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	kernel := internal.Kernel{
		WalletAddress: cfg.Wallet,
		ChainPort:     cfg.Ports.Http,
		// todo
	}

	return &kernel, nil
}
