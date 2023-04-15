package main

import (
	"bvpn-prototype/internal/cli/api"
	"bvpn-prototype/internal/cli/domain"
	"bvpn-prototype/internal/infrastructure/config"
	internal_di "bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/logger"
	"bvpn-prototype/internal/protocol/entity"
	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di"
	"gopkg.in/yaml.v3"
	"os"
)

const defaultConfigFile = "config.yaml"

var configFile = defaultConfigFile // todo: custom

func main() {
	command := os.Args[1]
	_, err := parseConfig()
	if err != nil {
		panic(err)
	}

	logger.Init()
	setupDi()

	cli := internal_di.Get("cli_api").(*api.CliApi)
	switch command {

	case "run":
	case "init":
		cli.Init()
		break

	case "make":
		switch os.Args[2] {

		case "transaction":
		case "tx":
			cli.MakeTx()
			break

		case "offer":
			cli.MakeOffer()
			break
		}
		break
	}
}

func setupDi() {
	app, _ := di.NewBuilder()

	app.Add(di.Def{
		Name: "cli_api",
		Build: func(ctn di.Container) (interface{}, error) {
			return api.NewCliApi()
		},
	})

	app.Add(di.Def{
		Name: "cli_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return domain.NewCliService()
		},
	})

	internal_di.Set(app.Build())
}

func parseConfig() (*config.Config, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var yamlCfg struct {
		Http struct {
			Port string `yaml:"port"`
			URL  string `yaml:"url" validate:"http_url"`
		} `yaml:"http"`
		VPN struct {
			Port  string `yaml:"port"`
			Proto string `yaml:"proto" validate:"oneof=udp tcp"`
		}
		Peers []struct {
			Ip      string `yaml:"ip" validate:"ip"`
			HttpUrl string `yaml:"url" validate:"http_url"`
		} `yaml:"peers" validate:"dive"`
	}

	err = yaml.Unmarshal(yamlFile, &yamlCfg)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(yamlCfg)
	if err != nil {
		return nil, err
	}

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
