package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNStoredTempTrade(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.StoredTempTrade {
	items := make([]types.StoredTempTrade, n)
	for i := range items {
		items[i].TradeIndex = uint64(i)

		keeper.SetStoredTempTrade(ctx, items[i])
	}
	return items
}

func TestStoredTempTradeGet(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	items := createNStoredTempTrade(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetStoredTempTrade(ctx,
			item.TradeIndex,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestStoredTempTradeRemove(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	items := createNStoredTempTrade(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveStoredTempTrade(ctx,
			item.TradeIndex,
		)
		_, found := keeper.GetStoredTempTrade(ctx,
			item.TradeIndex,
		)
		require.False(t, found)
	}
}

func TestStoredTempTradeGetAll(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	items := createNStoredTempTrade(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllStoredTempTrade(ctx)),
	)
}
