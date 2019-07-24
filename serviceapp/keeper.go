package serviceapp

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetService(ctx sdk.Context, hash []byte) []byte {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil
	}
	return store.Get(hash)
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetService(ctx sdk.Context, hash, data []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(hash, data)
}

func (k Keeper) GetServicesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
