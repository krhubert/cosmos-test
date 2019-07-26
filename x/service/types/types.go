package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("mesg", 1)}

// Whois is a struct that contains all the metadata of a name
type Service struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// Returns a new Whois with the minprice as the price
func NewService() Service {
	return Service{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (s Service) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, s.Owner, s.Value, s.Price))
}
