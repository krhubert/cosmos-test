package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
