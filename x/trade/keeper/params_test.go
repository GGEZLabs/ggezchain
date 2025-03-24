package keeper_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.TradeKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
