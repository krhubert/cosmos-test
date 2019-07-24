package service

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

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
