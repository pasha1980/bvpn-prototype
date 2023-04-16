package api

import (
	"bvpn-prototype/internal/cli/domain"
	"bvpn-prototype/internal/infrastructure/di"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"syscall"
)

type CliApi struct {
}

func (a *CliApi) Init() {
	var opts struct {
		Detached bool `short:"d" long:"detached" description:"Run in detached mode" required:"false"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		handle(err)
	}

	cliService := di.Get("cli_service").(*domain.CliService)
	err = cliService.Init()
	if err != nil {
		handle(err)
	}

	if !opts.Detached {
		ctlc := make(chan os.Signal)
		signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		<-ctlc
		close(ctlc)
	}
}

func (*CliApi) MakeTx() {
	var opts struct {
		To     string  `short:"t" long:"to" description:"Receiver" required:"true"`
		Amount float64 `short:"a" long:"amount" description:"Amount" required:"true"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		handle(err)
	}

	cliService := di.Get("cli_service").(*domain.CliService)
	err = cliService.MakeTx(opts.To, opts.Amount)
	if err != nil {
		handle(err)
	}
}

func (a *CliApi) MakeOffer() {
	var opts struct {
		Price float64 `short:"p" long:"price" description:"Price" required:"true"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		handle(err)
	}

	cliService := di.Get("cli_service").(*domain.CliService)
	err = cliService.MakeOffer(opts.Price)
	if err != nil {
		handle(err)
	}
}

func NewCliApi() (*CliApi, error) {
	return &CliApi{}, nil
}
