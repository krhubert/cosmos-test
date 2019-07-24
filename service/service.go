package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// Module ...
var Module = sdk.NewKVStoreKey("service")

// Cdc ...
var Cdc *amino.Codec

// Service ...
type Service struct {
	Hash []byte
	Name string
	Sid  string
	// ...
}

// RegisterCodec ...
func RegisterCodec(cdc *amino.Codec) {
	Cdc.RegisterConcrete(MsgSetService{}, Module.Name()+"/set_service", nil)
}
