package main

import (
	chain_domain "bvpn-prototype/internal/chain/domain"
	cli_api "bvpn-prototype/internal/cli/api"
	cli_domain "bvpn-prototype/internal/cli/domain"
	"bvpn-prototype/internal/infrastructure/config"
	internal_di "bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/logger"
	peer_domain "bvpn-prototype/internal/peer/domain"
	"bvpn-prototype/internal/protocol/entity"
	vpn_domain "bvpn-prototype/internal/vpn/domain"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di"
	"gopkg.in/yaml.v3"
	"os"
)

const defaultConfigFile = "./config.yaml"

var configFile = defaultConfigFile

func main() {
	logger.Init()
	setupDi()

	_, err := parseConfig()
	if err != nil {
		panic(err)
	}

	cli, _ := cli_api.NewCliApi()

	command := os.Args[1]
	switch command {
	case "init", "run":
		cli.Init()
		break

	case "make":
		fmt.Println(8)
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
		Name: "cli_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return cli_domain.NewCliService()
		},
	})

	app.Add(di.Def{
		Name: "chain_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return chain_domain.NewChainService()
		},
	})

	app.Add(di.Def{
		Name: "chain_public",
		Build: func(ctn di.Container) (interface{}, error) {
			return chain_domain.NewChainService()
		},
	})

	app.Add(di.Def{
		Name: "peer_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return peer_domain.NewPeerService()
		},
	})

	app.Add(di.Def{
		Name: "peer_public",
		Build: func(ctn di.Container) (interface{}, error) {
			return peer_domain.NewPeerService()
		},
	})

	app.Add(di.Def{
		Name: "vpn_public",
		Build: func(ctn di.Container) (interface{}, error) {
			return vpn_domain.NewVpnService()
		},
	})

	app.Add(di.Def{
		Name: "vpn_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return vpn_domain.NewVpnService()
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
			Port string `yaml:"port" validate:"number"`
			URL  string `yaml:"url" validate:"http_url"`
		} `yaml:"http"`
		VPN struct {
			Port  string `yaml:"port" validate:"number"`
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
