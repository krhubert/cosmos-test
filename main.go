package main

import (
	"os"

	"github.com/cosmos/sdk-application-tutorial/app"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	db := db.NewMemDB()
	app := app.NewServiceApp(logger, db)

	cfg := config.DefaultConfig()
	cfg.SetRoot("./node")
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}

	tmNode, err := node.NewNode(cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		app.Logger(),
	)
	if err != nil {
		panic(err)
	}
	if err := tmNode.Start(); err != nil {
		panic(err)
	}

	select {}
}
