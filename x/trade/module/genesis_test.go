package trade_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	trade "github.com/GGEZLabs/ggezchain/x/trade/module"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		TradeIndex: types.TradeIndex{
			NextId: 33,
		},
		StoredTradeList: []types.StoredTrade{
			{
				TradeIndex: 0,
			},
			{
				TradeIndex: 1,
			},
		},
		StoredTempTradeList: []types.StoredTempTrade{
			{
				TradeIndex: 0,
			},
			{
				TradeIndex: 1,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TradeKeeper(t)
	trade.InitGenesis(ctx, k, genesisState)
	got := trade.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.TradeIndex, got.TradeIndex)
	require.ElementsMatch(t, genesisState.StoredTradeList, got.StoredTradeList)
	require.ElementsMatch(t, genesisState.StoredTempTradeList, got.StoredTempTradeList)
	// this line is used by starport scaffolding # genesis/test/assert
}
