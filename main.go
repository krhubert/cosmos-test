package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/go-bip39"
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

const (
	moniker  = "cosmos-localhost"
	accName  = "bob"
	chainID  = "test-chain"
	password = "1234"
	denom    = "mesg"
)

func main() {
	// recreate directories for configs
	os.RemoveAll("config")
	os.RemoveAll("data")
	os.MkdirAll("config/gentx", 0755)
	os.MkdirAll("data", 0755)

	// create new default config
	cfg := config.DefaultConfig()
	cfg.Moniker = moniker
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
	info, err := kb.CreateAccount(accName, mnemonic, "", password, 0, 0)
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

	genAcc := genaccounts.NewGenesisAccountRaw(addr, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(10000))), sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(10000))), 0, 1, "")
	if err := genAcc.Validate(); err != nil {
		panic(err)
	}

	appState := simapp.ModuleBasics.DefaultGenesis()

	// add genesis account to the app state
	genesisAccounts := []genaccounts.GenesisAccount{genAcc}
	appState[genaccounts.ModuleName] = cdc.MustMarshalJSON(genaccounts.GenesisState(genesisAccounts))

	if err = simapp.ModuleBasics.ValidateGenesis(appState); err != nil {
		panic(err)
	}

	fmt.Println(string(codec.MustMarshalJSONIndent(cdc, appState)))

	// generate first tx
	nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(cfg) // get node id and pub key
	if err != nil {
		panic(err)
	}

	// create msg for validato creation
	description := staking.NewDescription(moniker, "", "", "")
	minSelfDelegation := sdk.NewInt(1)
	selfDelegationCoin := sdk.NewCoin(denom, sdk.NewInt(1))
	msg := staking.NewMsgCreateValidator(sdk.ValAddress(addr), info.GetPubKey(), selfDelegationCoin, description, staking.CommissionRates{}, minSelfDelegation)

	encoder := auth.DefaultTxEncoder(cdc)
	txBldr := auth.NewTxBuilder(encoder, 0, 0, 0, 0, false, chainID, mnemonic, nil, nil)

	// buidl standard sign msg with txbuilder and calculate the fees
	stdSignMsg, err := txBldr.BuildSignMsg([]sdk.Msg{msg})
	if err != nil {
		panic(err)
	}

	// create std tx
	stdTx, err := auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo), nil
	if err != nil {
		panic(err)
	}

	// sign the tx
	signTx, err := txBldr.SignStdTx(accName, password, stdTx, false)
	if err != nil {
		panic(err)
	}

	// write first sign tx to file
	genTxFile := filepath.Join(genTxsDir, fmt.Sprintf("gentx-%v.json", nodeID))
	if err := ioutil.WriteFile(genTxFile, cdc.MustMarshalJSON(signTx), 0644); err != nil {
		panic(err)
	}

	// set genesic docs and export to file
	genDoc := &ttypes.GenesisDoc{}
	genDoc.ChainID = chainID
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
	fmt.Println(">>>>>> app-message <<<<<<")
	fmt.Println(string(appMessage))
	fmt.Println(">>>>>> app-message <<<<<<")

	// generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}

	// start sapp as standalone server for test.
	// s, err := server.NewServer(cfg.ProxyApp, cfg.ABCI, sapp)
	// if err != nil {
	// 	panic(err)
	// }
	// s.SetLogger(logger)
	// if err := s.Start(); err != nil {
	// 	panic(err)
	// }

	fmt.Println(">>>>>> create node <<<<<<")
	node, err := node.NewNode(cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(sapp), // sapp could be also started as local client
		// proxy.NewRemoteClientCreator(cfg.ProxyApp, cfg.ABCI, false),
		// proxy.NewLocalClientCreator(counter.NewCounterApplication(false)),
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
