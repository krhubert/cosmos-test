package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgPublish ...
type MsgPublish struct {
	Sid   string         `json:"sid"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgPublish ...
func NewMsgPublish(sid string, owner sdk.AccAddress) MsgPublish {
	return MsgPublish{
		Sid:   sid,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgPublish) Route() string { return ModuleName }

// Type should return the action
func (msg MsgPublish) Type() string { return "publish" }

// ValidateBasic runs stateless checks on the message
func (msg MsgPublish) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if msg.Sid == "" {
		return sdk.ErrUnknownRequest("sid cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgPublish) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgPublish) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
