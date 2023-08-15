package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

func CmdListStoredTrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-stored-trade",
		Short: "list all storedTrade",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllStoredTradeRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.StoredTradeAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowStoredTrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-stored-trade [trade-index]",
		Short: "shows a storedTrade",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argTradeIndex, err := cast.ToUint64E(args[0])

			params := &types.QueryGetStoredTradeRequest{
				TradeIndex: argTradeIndex,
			}

			res, err := queryClient.StoredTrade(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
