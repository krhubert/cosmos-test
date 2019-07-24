package serviceapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Service struct {
	Hash []byte
	Data []byte
}

type GenesisState struct {
	Services []*Service `json:"services"`
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(gs GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Services: make([]*Service, 0),
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, gs GenesisState) []abci.ValidatorUpdate {
	for _, s := range gs.Services {
		keeper.SetService(ctx, s.Hash, s.Data)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var services []*Service

	iterator := k.GetServicesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		hash := iterator.Key()
		data := k.GetService(ctx, hash)
		services = append(services, &Service{Hash: hash, Data: data})
	}
	return GenesisState{Services: services}
}
