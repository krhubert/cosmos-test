package service

import (
	"github.com/cosmos/sdk-application-tutorial/x/service/types"
)

// ...
const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

// ...
var (
	NewMsgBuyName = types.NewMsgBuyName
	NewMsgSetName = types.NewMsgSetName
	NewService    = types.NewService
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	// MsgSetName ...
	MsgSetName = types.MsgSetName
	// MsgBuyName ...
	MsgBuyName = types.MsgBuyName
	// QueryResResolve ...
	QueryResResolve = types.QueryResResolve
	// QueryResNames ...
	QueryResNames = types.QueryResNames
	// Service ...
	Service = types.Service
)
