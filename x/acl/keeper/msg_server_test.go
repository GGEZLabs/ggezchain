package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, context.Context) {
	k, ctx := keepertest.AclKeeper(t)
	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}
