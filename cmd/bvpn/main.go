package main

import (
	"bvpn-prototype/internal/api"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/storage/config"
	"fmt"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const defaultConfigFile = "config.yaml"

var configFile = defaultConfigFile

func main() {
	var opts struct {
		ConfigFile *string  `short:"c" long:"configs" description:"Configuration file" required:"false"`
		To         *string  `short:"t" long:"to" description:"Receiver" required:"false"`
		Amount     *float64 `short:"a" long:"amount" description:"Amount" required:"false"`
		Price      *float64 `short:"p" long:"price" description:"Price" required:"false"`
		Detached   bool     `short:"d" long:"detached" description:"Run in detached mode" required:"false"`
	}

	commands, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err)
	}
	command := commands[0]

	if opts.ConfigFile != nil {
		configFile = *opts.ConfigFile
	}
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalln(err)
	}

	cliApi := api.CLI(cfg)

	switch command {
	case "run":
	case "init":
		cliApi.Init(opts.Detached)
		break
	case "make":
		switch commands[1] {
		case "transaction":
		case "tx":

			if opts.To == nil {
				log.Fatalln("To whom?")
			}

			if opts.Amount == nil {
				log.Fatalln("Should I send everything you have????")
			}

			cliApi.MakeTx(
				*opts.To,
				*opts.Amount,
			)
			break
		case "offer":
			if opts.Price == nil {
				log.Fatalln("I'll sell you to devil, if you will not tell me the price!")
			}

			cliApi.MakeOffer(*opts.Price)

			break
		}
		break
	default:
		fmt.Println("Hello") // todo
	}
}

func parseConfig() (*config.Config, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var yamlCfg struct {
		Http struct {
			Port string `yaml:"port"`
			URL  string `yaml:"url"`
		} `yaml:"http"`
		VPN struct {
			Port  string `yaml:"port"`
			Proto string `yaml:"proto"`
		}
		Peers []struct {
			Ip      string `yaml:"ip"`
			HttpUrl string `yaml:"url"`
		} `yaml:"peers"`
	}

	err = yaml.Unmarshal(yamlFile, &yamlCfg)
	if err != nil {
		return nil, err
	}

	// todo: validation

	var peers []entity.Node
	for _, peerCfg := range yamlCfg.Peers {
		peers = append(peers, entity.Node{
			URL: peerCfg.HttpUrl,
			IP:  peerCfg.Ip,
		})
	}

	cfg := config.Config{
		StorageDirectory: ".",
		URL:              yamlCfg.Http.URL,
		HttpPort:         yamlCfg.Http.Port,
		VpnPort:          yamlCfg.VPN.Port,
		VpnProto:         yamlCfg.VPN.Proto,
		Peers:            peers,
	}

	config.Set(cfg)

	return &cfg, nil
}
