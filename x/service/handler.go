package service

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "service" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPublish:
			return handleMsgPublish(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgPublish(ctx sdk.Context, keeper Keeper, msg MsgPublish) sdk.Result {
	service := keeper.GetService(ctx, msg.Sid)
	if !service.Owner.Empty() {
		return sdk.ErrUnauthorized("Already deployed").Result()
	}
	keeper.SetService(ctx, msg.Sid, Service{msg.Sid, msg.Owner})
	return sdk.Result{}
}
