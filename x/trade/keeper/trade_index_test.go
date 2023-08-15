package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

func createTestTradeIndex(keeper *keeper.Keeper, ctx sdk.Context) types.TradeIndex {
	item := types.TradeIndex{}
	keeper.SetTradeIndex(ctx, item)
	return item
}

func TestTradeIndexGet(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	item := createTestTradeIndex(keeper, ctx)
	rst, found := keeper.GetTradeIndex(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestTradeIndexRemove(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	createTestTradeIndex(keeper, ctx)
	keeper.RemoveTradeIndex(ctx)
	_, found := keeper.GetTradeIndex(ctx)
	require.False(t, found)
}
