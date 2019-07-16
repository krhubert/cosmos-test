package main

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func main() {
	// create new default config
	cfg := config.DefaultConfig()
	cfg.SetRoot(simapp.DefaultNodeHome)

	// create new json logger
	logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))

	// create database for app
	db := db.NewMemDB()

	// create simple app
	sapp := simapp.NewSimApp(logger, db, nil, true, 0)

	// generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}

	fmt.Println(">>>>>> create node <<<<<<")
	node, err := node.NewNode(cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(sapp),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(">>>>>> start the server <<<<<<")
	if err := node.Start(); err != nil {
		panic(err)
	}

	select {}
}
