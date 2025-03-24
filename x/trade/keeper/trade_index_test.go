package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"
)

func createTestTradeIndex(keeper keeper.Keeper, ctx context.Context) types.TradeIndex {
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
