package query

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/sdk-application-tutorial/x/service/types"
	"github.com/spf13/cobra"
)



func query(ctx context.CLIContext, cdc *codec.Codec, path string, out interface{}) error {
	res, _, err := ctx.QueryWithData(path, nil)
	if err != nil {
		return err
	}
	cdc.MustUnmarshalJSON(res, &out)
	return nil
}

// Resolve ...
func Resolve(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "resolve [name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)
			var out types.QueryResResolve
			if err := query(ctx, cdc, fmt.Sprintf("custom/service/resolve/%s", args[0]), out); err != nil {
				return err
			}
			return ctx.PrintOutput(out)
		},
	}
}

// GetService ...
func GetService(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "service [name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)
			var out types.QueryResResolve
			if err := query(ctx, cdc, fmt.Sprintf("custom/service/resolve/%s", args[0]), out); err != nil {
				return err
			}
			return ctx.PrintOutput(out)
		},
	}
}

// GetNames ...
func GetNames(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "names",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)
			var out types.QueryResResolve
			if err := query(ctx, cdc, "custom/service/names", out); err != nil {
				return err
			}
			return ctx.PrintOutput(out)
		},
	}
}
