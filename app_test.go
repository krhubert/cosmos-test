package main

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestSimAppExport(t *testing.T) {
	db := db.NewMemDB()
	app, err := NewSimApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	require.NoError(t, err)

	genesisState := ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2, err := NewSimApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	require.NoError(t, err)

	_, _, err = app2.ExportAppStateAndValidators()
	require.NoError(t, err)
}
