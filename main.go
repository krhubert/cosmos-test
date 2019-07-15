package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	ctypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	ttypes "github.com/tendermint/tendermint/types"
	// "github.com/tendermint/tendermint/p2p"
	// "github.com/tendermint/tendermint/privval"
	// "github.com/tendermint/tendermint/proxy"
)

const accName = "bob"

func main() {
	// recreate directories for configs
	os.RemoveAll("config")
	os.RemoveAll("data")
	os.MkdirAll("config/gentx", 0755)
	os.MkdirAll("data", 0755)

	// create new default config
	cfg := config.DefaultConfig()
	genTxsDir := filepath.Join(cfg.RootDir, "config", "gentx")

	// create new json logger
	logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))

	// create database for app
	db := db.NewMemDB()

	// get and init codec
	cdc := simapp.MakeCodec()

	// create simple app
	sapp := simapp.NewSimApp(logger, db, nil, true, 0)

	// create key base for accounts - could be in memory (keys.NewInMem())
	kb, err := keys.NewKeyBaseFromDir(cfg.RootDir)
	if err != nil {
		panic(err)
	}
	// read entropy seed straight from crypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		panic(err)
	}

	// use bip39 for mnemonic generation
	mnemonic, err := bip39.NewMnemonic(entropySeed[:])
	if err != nil {
		panic(err)
	}

	// create actual account
	info, err := kb.CreateAccount(accName, mnemonic, "", "", 0, 0)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "public key:", info.GetPubKey())
	fmt.Fprintln(os.Stderr, "address:", info.GetAddress())
	fmt.Fprintln(os.Stderr, "\n**Important** write this mnemonic phrase in a safe place.")
	fmt.Fprintln(os.Stderr, "It is the only way to recover your account if you ever forget your password.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, mnemonic)
	fmt.Fprintln(os.Stderr, "")

	addr := info.GetAddress()

	genAcc := genaccounts.NewGenesisAccountRaw(addr, ctypes.NewCoins(), ctypes.NewCoins(), 0, 0, "")
	if err := genAcc.Validate(); err != nil {
		panic(err)
	}

	appState := simapp.ModuleBasics.DefaultGenesis()

	// add genesis account to the app state
	genesisAccounts := []genaccounts.GenesisAccount{genAcc}

	appState[genaccounts.ModuleName] = cdc.MustMarshalJSON(genaccounts.GenesisState(genesisAccounts))

	fmt.Println(string(codec.MustMarshalJSONIndent(cdc, appState)))

	// generate first tx
	if err = simapp.ModuleBasics.ValidateGenesis(appState); err != nil {
		panic(err)
	}

	nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(cfg)
	if err != nil {
		panic(err)
	}

	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := client.NewCLIContext().WithCodec(cdc)

	tx, err := utils.SignStdTx(txBldr, cliCtx, name, stdTx, false, true)
	if err != nil {
		panic(err)
	}

	genTxFile := filepath.Join(genTxsDir, fmt.Sprintf("gentx-%v.json", nodeID))
	if err := ioutil.WriteFile(genTxFile, cdc.MustMarshalJSON(tx), 0644); err != nil {
		panic(err)
	}

	// set genesic docs and export to file
	genDoc := &ttypes.GenesisDoc{}
	genDoc.ChainID = "SimApp"
	genDoc.AppState = codec.MustMarshalJSONIndent(cdc, appState)
	if err := genutil.ExportGenesisFile(genDoc, cfg.GenesisFile()); err != nil {
		panic(err)
	}

	// collect-gentxs

	initCfg := genutil.NewInitConfig(genDoc.ChainID, genTxsDir, "", nodeID, valPubKey)

	appMessage, err := genutil.GenAppStateFromConfig(cdc, cfg, initCfg, *genDoc, genaccounts.AppModuleBasic{})
	if err != nil {
		panic(err)
	}
	fmt.Println(appMessage)

	// generate node PrivKey
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
		// proxy.NewLocalClientCreator(counter.NewCounterApplication(false)),
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

// initialize the app chain
// genesisState := simapp.NewDefaultGenesisState()
// stateBytes, err := codec.MarshalJSONIndent(cdc, genesisState)
// if err != nil {
// 	panic(err)
// }

// sapp.InitChain(
// 	abci.RequestInitChain{
// 		Validators:    []abci.ValidatorUpdate{},
// 		AppStateBytes: stateBytes,
// 	},
// )
// sapp.Commit()

// appState, validators, err := sapp.ExportAppStateAndValidators(false, []string{})
// if err != nil {
// 	panic(err)
// }
