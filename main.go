package main

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/test/app"
	"github.com/mesg-foundation/test/service"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	tm "github.com/tendermint/tendermint/types"
)

func main() {
	app := app.New()

	cfg := config.DefaultConfig()
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}

	validator := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
	node, err := node.NewNode(cfg,
		validator,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genesisLoader(nodeKey),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		app.Logger(),
	)
	if err != nil {
		panic(err)
	}
	if err := node.Start(); err != nil {
		panic(err)
	}

	time.Sleep(1000)

	createService(app)

	select {}
}

func createService(application *app.MyApp) sdk.Result {
	return application.Deliver(app.NewTx([]sdk.Msg{
		service.MsgSetService{
			Name: "test",
			Sid:  "test",
		},
	}))
}

func genesisLoader(validator *p2p.NodeKey) func() (*tm.GenesisDoc, error) {
	return func() (*tm.GenesisDoc, error) {
		return &tm.GenesisDoc{
			GenesisTime:     time.Unix(0, 0),
			ChainID:         "xxx",
			ConsensusParams: tm.DefaultConsensusParams(),
			Validators: []tm.GenesisValidator{
				tm.GenesisValidator{
					Address: validator.PubKey().Address(),
					PubKey:  validator.PubKey(),
					Power:   1,
					Name:    "validator",
				},
			},
		}, nil
	}
}
