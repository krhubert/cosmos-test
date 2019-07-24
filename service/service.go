package service

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Module ...
var Module = sdk.NewKVStoreKey("service")

// Cdc ...
var Cdc *amino.Codec

/*******
* DATA *
*******/

// Service ...
type Service struct {
	Hash []byte
	Name string
	Sid  string
	// ...
}

/*********
* KEEPER *
*********/

// Keeper ...
type Keeper struct {
	cdc *codec.Codec
}

// NewKeeper ...
func NewKeeper(cdc *codec.Codec) Keeper {
	return Keeper{cdc: cdc}
}

// Find ...
func (k Keeper) Find(ctx sdk.Context, hash []byte) (Service, error) {
	store := ctx.KVStore(Module)
	if !store.Has(hash) {
		return Service{}, fmt.Errorf("not found")
	}
	s := store.Get(hash)
	var service Service
	return service, k.cdc.UnmarshalBinaryBare(s, &service)
}

// Add ...
func (k Keeper) Add(ctx sdk.Context, hash []byte, service Service) error {
	store := ctx.KVStore(Module)
	value, err := k.cdc.MarshalBinaryBare(service)
	if err != nil {
		return err
	}
	store.Set(hash, value)
	return nil
}

/**********
* QUERIER *
**********/

// NewQuerier ...
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case "find":
			service, err := keeper.Find(ctx, []byte(path[1]))
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			res, err := codec.MarshalJSONIndent(keeper.cdc, service)
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			return res, nil
		default:
			return nil, sdk.ErrUnknownRequest("invalid path")
		}
	}
}

/******
* MSG *
*******/

// MsgSetService ...
type MsgSetService struct {
	Name string
	Sid  string
}

// Route ...
func (msg MsgSetService) Route() string { return Module.Name() }

// Type ...
func (msg MsgSetService) Type() string { return "set_service" }

// ValidateBasic ...
func (msg MsgSetService) ValidateBasic() sdk.Error { return nil }

// GetSignBytes ...
func (msg MsgSetService) GetSignBytes() []byte { return sdk.MustSortJSON(Cdc.MustMarshalJSON(msg)) }

// GetSigners ...
func (msg MsgSetService) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{} }

/**********
* HANDLER *
**********/

// NewHandler ...
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetService:
			hash := []byte("calculate_hash_of_service")
			if err := keeper.Add(ctx, hash, Service{
				Hash: hash,
				Name: msg.Name,
				Sid:  msg.Sid,
			}); err != nil {
				return sdk.ErrInsufficientCoins(err.Error()).Result()
			}
			return sdk.Result{
				Data: hash,
			}
		default:
			return sdk.ErrUnknownRequest("unknown msg").Result()
		}
	}
}

// RegisterCodec ...
func RegisterCodec(cdc *amino.Codec) {
	Cdc.RegisterConcrete(MsgSetService{}, Module.Name()+"/set_service", nil)
}
