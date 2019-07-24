package app

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/mesg-foundation/test/service"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

var mainStoreKey = sdk.NewKVStoreKey(baseapp.MainStoreKey)

// MyApp ...
type MyApp struct {
	*baseapp.BaseApp
}

// New ...
func New() *MyApp {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	app := &MyApp{
		BaseApp: baseapp.NewBaseApp("test", log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db.NewMemDB(), auth.DefaultTxDecoder(cdc)),
	}

	keeper := service.NewKeeper(cdc)
	fmt.Println("keeper")
	app.Router().AddRoute(service.Module.Name(), service.NewHandler(keeper))
	fmt.Println("handler created")
	app.QueryRouter().AddRoute(service.Module.Name(), service.NewQuerier(keeper))
	fmt.Println("querier created")
	app.MountStores(mainStoreKey)
	app.LoadLatestVersion(mainStoreKey)

	return app
}

// Tx ...
type Tx struct {
	msgs []sdk.Msg
}

// NewTx ...
func NewTx(messages []sdk.Msg) Tx {
	return Tx{messages}
}

// GetMsgs ...
func (t Tx) GetMsgs() []sdk.Msg {
	return t.msgs
}

// ValidateBasic ...
func (t Tx) ValidateBasic() sdk.Error {
	for _, msg := range t.msgs {
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}
