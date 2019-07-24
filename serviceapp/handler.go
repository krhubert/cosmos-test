package serviceapp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetService:
			keeper.SetService(ctx, msg.Hash, msg.Data)
			return sdk.Result{}
		case MsgGetService:
			return sdk.Result{Data: keeper.GetService(ctx, msg.Hash)}
		default:
			errMsg := fmt.Sprintf("Unrecognized serviceapp Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
