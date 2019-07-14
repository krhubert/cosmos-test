package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/tendermint/tendermint/abci/server"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
	// "github.com/tendermint/tendermint/p2p"
	// "github.com/tendermint/tendermint/privval"
	// "github.com/tendermint/tendermint/proxy"
)

func main() {
	// recreate directories for configs
	os.RemoveAll("config")
	os.RemoveAll("data")
	os.MkdirAll("config", 0755)
	os.MkdirAll("data", 0755)

	// create new default config
	cfg := config.DefaultConfig()

	// create new json logger
	logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))

	// create database for app
	db := db.NewMemDB()

	// get and init codec
	cdc := simapp.MakeCodec()

	// create simple app
	sapp := simapp.NewSimApp(logger, db, nil, true, 0)

	// initialize the app chain
	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(cdc, genesisState)
	if err != nil {
		panic(err)
	}

	sapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	sapp.Commit()

	appState, validators, err := sapp.ExportAppStateAndValidators(false, []string{})
	if err != nil {
		panic(err)
	}

	// get only app state without validators (init chain is not required)
	// appState, err := codec.MarshalJSONIndent(cdc, simapp.ModuleBasics.DefaultGenesis())
	// if err != nil {
	// 	panic(err)
	// }

	genDoc := &types.GenesisDoc{}
	genDoc.ChainID = "SimApp"
	genDoc.Validators = validators
	genDoc.AppState = appState
	if err := genutil.ExportGenesisFile(genDoc, cfg.GenesisFile()); err != nil {
		panic(err)
	}

	// Generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}

	// start sapp as standalone server for test.
	s, err := server.NewServer(cfg.ProxyApp, cfg.ABCI, sapp)
	if err != nil {
		panic(err)
	}
	s.SetLogger(logger)
	if err := s.Start(); err != nil {
		panic(err)
	}

	node, err := node.NewNode(cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		// proxy.NewLocalClientCreator(sapp), // sapp could be also started as local client
		proxy.NewRemoteClientCreator(cfg.ProxyApp, cfg.ABCI, false),
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
}
