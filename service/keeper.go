package service

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
