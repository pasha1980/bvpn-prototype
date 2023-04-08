package main

import (
	"bvpn-prototype/internal"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/storage/config"
	"fmt"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const defaultConfigFile = "config.yaml"

var configFile = defaultConfigFile

func main() {
	ctlc := make(chan os.Signal)
	signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	var opts struct {
		ConfigFile *string  `short:"c" long:"configs" description:"Configuration file" required:"false"`
		To         *string  `short:"t" long:"to" description:"Receiver" required:"false"`
		Amount     *float64 `short:"a" long:"amount" description:"Amount" required:"false"`
		Price      *float64 `short:"p" long:"price" description:"Price" required:"false"`
	}

	commands, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err)
	}
	command := commands[0]

	if opts.ConfigFile != nil {
		configFile = *opts.ConfigFile
	}
	kernel, err := createKernel()
	if err != nil {
		log.Fatalln(err)
	}

	switch command {
	case "run":
		kernel.Run()
		<-ctlc
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

			kernel.MakeTx(
				*opts.To,
				*opts.Amount,
			)
			break
		case "offer":
			if opts.Price == nil {
				log.Fatalln("I'll sell you to devil, if you will not tell me the price!")
			}

			kernel.MakeOffer(*opts.Price)

			break
		}
		break
	default:
		fmt.Println("Hello") // todo
	}

	close(ctlc)
}

func createKernel() (*internal.CliApi, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg struct {
		Url   string `yaml:"url"`
		Port  uint64 `yaml:"port"`
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

	if cfg.Port == 0 {
		cfg.Port = 80
	}

	kernel := internal.CliApi{
		URL:      cfg.Url,
		HttpPort: cfg.Port,
		Peers:    peers,
	}

	config.Set(config.Config{
		StorageDirectory: ".",
		URL:              cfg.Url,
	})

	return &kernel, nil
}
