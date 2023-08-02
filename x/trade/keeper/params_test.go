package keeper_test

import (
	"testing"

	testkeeper "github.com/GGEZLabs/testchain/testutil/keeper"
	"github.com/GGEZLabs/testchain/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TradeKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
