package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
