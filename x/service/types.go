package service

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Service ...
type Service struct {
	Sid   string         `json:"sid"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewService ...
func NewService() Service {
	return Service{}
}

// implement fmt.Stringer
func (s Service) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s, Sid: %s`, s.Owner, s.Sid))
}
