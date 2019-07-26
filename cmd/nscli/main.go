package main

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	app "github.com/cosmos/sdk-application-tutorial/app"
	"github.com/cosmos/sdk-application-tutorial/cmd/nscli/query"
	"github.com/cosmos/sdk-application-tutorial/cmd/nscli/tx"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

func main() {
	cdc := app.MakeCodec()
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(
		queryCmd(cdc),
		txCmd(cdc),
		keys.Commands(),
	)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use: "query",
	}
	queryCmd.AddCommand(
		query.Resolve(cdc),
		query.GetService(cdc),
		query.GetNames(cdc),
	)
	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use: "tx",
	}
	txCmd.AddCommand(
		tx.BuyName(cdc),
		tx.SetName(cdc),
	)
	return txCmd
}
