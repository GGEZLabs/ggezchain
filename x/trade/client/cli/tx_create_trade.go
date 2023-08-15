package cli

import (
	"strconv"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateTrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-trade [trade-type] [coin] [price] [quantity] [receiver-address] [trade-data]",
		Short: "Broadcast message createTrade",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTradeType := args[0]
			argCoin := args[1]
			argPrice := args[2]
			argQuantity := args[3]
			argReceiverAddress := args[4]
			argTradeData := args[5]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateTrade(
				clientCtx.GetFromAddress().String(),
				argTradeType,
				argCoin,
				argPrice,
				argQuantity,
				argReceiverAddress,
				argTradeData,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
