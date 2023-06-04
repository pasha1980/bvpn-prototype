package main

import (
	chain_gateway "bvpn-prototype/internal/chain/api_gateway"
	chain_domain "bvpn-prototype/internal/chain/domain"
	chain_storage "bvpn-prototype/internal/chain/storage"
	cli_api "bvpn-prototype/internal/cli/api"
	cli_domain "bvpn-prototype/internal/cli/domain"
	"bvpn-prototype/internal/infrastructure/config"
	internal_di "bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/logger"
	peer_gateway "bvpn-prototype/internal/peer/api_gateway"
	peer_domain "bvpn-prototype/internal/peer/domain"
	peer_storage "bvpn-prototype/internal/peer/storage"
	"bvpn-prototype/internal/protocol/entity"
	vpn_domain "bvpn-prototype/internal/vpn/domain"
	vpn_storage "bvpn-prototype/internal/vpn/storage"
	"bvpn-prototype/tests"
	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di"
	"gopkg.in/yaml.v3"
	"os"
)

const defaultConfigFile = "./config.yaml"

var configFile = defaultConfigFile

func main() {
	cli, _ := cli_api.NewCliApi()

	if len(os.Args) == 1 {
		cli.Introduce()
		return
	}

	command := os.Args[1]
	switch command {
	case "test", "check":
		setupDi("test")
		cli.Test()
		break

	case "init", "run":
		preparation()
		cli.Init()
		break

	case "make":
		preparation()
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

func preparation() {
	checkMe()
	logger.Init()
	setupDi("live")

	_, err := parseConfig()
	if err != nil {
		panic(err)
	}
}

func checkMe() {
	result := tests.Run(tests.TypeInfrastructure)
	if !result.IsSucceed() {
		result.Print()
		os.Exit(1)
	}
}

func setupDi(env string) {
	app, _ := di.NewBuilder()

	// todo: test dependencies

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

	app.Add(di.Def{
		Name: "chain_repo",
		Build: func(ctn di.Container) (interface{}, error) {
			return chain_storage.NewChainRepo()
		},
	})

	app.Add(di.Def{
		Name: "mempool",
		Build: func(ctn di.Container) (interface{}, error) {
			return chain_storage.NewMempoolRepo()
		},
	})

	app.Add(di.Def{
		Name: "peer_repo",
		Build: func(ctn di.Container) (interface{}, error) {
			return peer_storage.NewPeerRepo()
		},
	})

	app.Add(di.Def{
		Name: "vpn_profile_repo",
		Build: func(ctn di.Container) (interface{}, error) {
			return vpn_storage.NewProfileRepo()
		},
	})

	app.Add(di.Def{
		Name: "chain_api_gateway",
		Build: func(ctn di.Container) (interface{}, error) {
			return chain_gateway.NewChainApiGateway()
		},
	})

	app.Add(di.Def{
		Name: "peer_api_gateway",
		Build: func(ctn di.Container) (interface{}, error) {
			return peer_gateway.NewPeerApiGateway()
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
