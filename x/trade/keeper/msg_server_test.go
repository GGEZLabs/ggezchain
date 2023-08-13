package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/GGEZLabs/testchain/testutil/keeper"
	"github.com/GGEZLabs/testchain/x/trade/keeper"
	"github.com/GGEZLabs/testchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t *testing.T) (types.MsgServer, context.Context) {
	k, ctx := keepertest.TradeKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
