package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"github.com/GGEZLabs/ggezchain/testutil/network"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/trade/client/cli"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

func networkWithTradeIndexObjects(t *testing.T) (*network.Network, types.TradeIndex) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	tradeIndex := types.TradeIndex{}
	nullify.Fill(&tradeIndex)
	state.TradeIndex = tradeIndex
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.TradeIndex
}

func TestShowTradeIndex(t *testing.T) {
	net, obj := networkWithTradeIndexObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc string
		args []string
		err  error
		obj  types.TradeIndex
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowTradeIndex(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetTradeIndexResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.TradeIndex)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.TradeIndex),
				)
			}
		})
	}
}
