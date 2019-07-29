package service

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case "find":
			srv := keeper.GetService(ctx, path[0])
			if srv.Sid == "" {
				return nil, sdk.ErrInternal("service not found")
			}
			res, err := codec.MarshalJSONIndent(keeper.cdc, srv)
			if err != nil {
				return nil, sdk.ErrInternal("could not marshal result to JSON")
			}

			return res, nil
		default:
			return nil, sdk.ErrUnknownRequest("unknown service query endpoint")
		}
	}
}
