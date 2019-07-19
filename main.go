package main

import (
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
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
	cfg.SetRoot("s1")
	// cfg.P2P.Seeds = "57e1d834e500ac497cc30c884ebb85a0879ca295@localhost:26656"
	// cfg.RPC.ListenAddress = "tcp://0.0.0.0:36657"
	// cfg.P2P.ListenAddress = "tcp://0.0.0.0:36656"

	// create new json logger
	// logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))
	logger := log.NewNopLogger()

	// create database for app
	db := db.NewMemDB()

	// create simple app
	sapp, err := NewSimApp(logger, db)
	if err != nil {
		panic(err)
	}

	// generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}

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

	if err := node.Start(); err != nil {
		panic(err)
	}

	fmt.Println("LET DO THIS")
	time.Sleep(1 * time.Second)
	ctx := sapp.NewContext(true, abci.Header{})
	sapp.saKeeper.SetService(ctx, []byte("test"), []byte("my service"))
	fmt.Println(sapp.saKeeper.GetService(ctx, []byte("test")))

	select {}
}
