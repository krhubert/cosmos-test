package serviceapp

import sdk "github.com/cosmos/cosmos-sdk/types"

const RouterKey = "serviceapp"

// MsgSetService defines a SetService message
type MsgSetService struct {
	Hash []byte `json:"hash"`
	Data []byte `json:"data"`
}

// NewMsgSetService is a constructor function for MsgSetService
func NewMsgSetService(hash, data []byte) MsgSetService {
	return MsgSetService{
		Hash: hash,
		Data: data,
	}
}

// Route should return the hash of the module
func (msg MsgSetService) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetService) Type() string { return "set_service" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetService) ValidateBasic() sdk.Error {
	if len(msg.Hash) == 0 || len(msg.Data) == 0 {
		return sdk.ErrUnknownRequest("Hash and/or Data cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

// MsgGetService defines a GetService message
type MsgGetService struct {
	Hash []byte `json:"hash"`
}

// NewMsgGetService is a constructor function for MsgGetService
func NewMsgGetService(hash []byte) MsgGetService {
	return MsgGetService{
		Hash: hash,
	}
}

// Route should return the hash of the module
func (msg MsgGetService) Route() string { return RouterKey }

// Type should return the action
func (msg MsgGetService) Type() string { return "get_service" }

// ValidateBasic runs stateless checks on the message
func (msg MsgGetService) ValidateBasic() sdk.Error {
	if len(msg.Hash) == 0 {
		return sdk.ErrUnknownRequest("Hash cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgGetService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgGetService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
