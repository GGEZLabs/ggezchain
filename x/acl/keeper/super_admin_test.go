package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func createTestSuperAdmin(keeper keeper.Keeper, ctx context.Context) types.SuperAdmin {
	item := types.SuperAdmin{
		Admin: sample.AccAddress(),
	}
	keeper.SetSuperAdmin(ctx, item)
	return item
}

func TestSuperAdminGet(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	item := createTestSuperAdmin(keeper, ctx)
	rst, found := keeper.GetSuperAdmin(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestSuperAdminRemove(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	createTestSuperAdmin(keeper, ctx)
	keeper.RemoveSuperAdmin(ctx)
	_, found := keeper.GetSuperAdmin(ctx)
	require.False(t, found)
}
