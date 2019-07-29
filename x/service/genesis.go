package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState ...
type GenesisState struct {
	Services []Service `json:"services"`
}

// NewGenesisState ...
func NewGenesisState(services []Service) GenesisState {
	return GenesisState{Services: nil}
}

// ValidateGenesis ...
func ValidateGenesis(data GenesisState) error {
	return nil
}

// DefaultGenesisState ...
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Services: []Service{},
	}
}

// InitGenesis ...
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.Services {
		keeper.SetService(ctx, record.Sid, record)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Service
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		var service Service
		service = k.GetService(ctx, name)
		records = append(records, service)
	}
	return GenesisState{Services: records}
}
