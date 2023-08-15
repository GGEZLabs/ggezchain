package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

func CmdListStoredTempTrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-stored-temp-trade",
		Short: "list all storedTempTrade",
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

			params := &types.QueryAllStoredTempTradeRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.StoredTempTradeAll(cmd.Context(), params)
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

func CmdShowStoredTempTrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-stored-temp-trade [trade-index]",
		Short: "shows a storedTempTrade",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argTradeIndex, err := cast.ToUint64E(args[0])

			params := &types.QueryGetStoredTempTradeRequest{
				TradeIndex: argTradeIndex,
			}

			res, err := queryClient.StoredTempTrade(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
