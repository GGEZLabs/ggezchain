package keeper_test

import (
	"context"
	"strconv"
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNAclAdmin(keeper keeper.Keeper, ctx context.Context, n int) []types.AclAdmin {
	items := make([]types.AclAdmin, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)

		keeper.SetAclAdmin(ctx, items[i])
	}
	return items
}

func TestSetAclAdmins(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := make([]types.AclAdmin, 10)
	keeper.SetAclAdmins(ctx, items)
	for _, item := range items {
		rst, found := keeper.GetAclAdmin(ctx,
			item.Address,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestAclAdminGet(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := createNAclAdmin(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetAclAdmin(ctx,
			item.Address,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestAclAdminRemove(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := createNAclAdmin(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAclAdmin(ctx,
			item.Address,
		)
		_, found := keeper.GetAclAdmin(ctx,
			item.Address,
		)
		require.False(t, found)
	}
}

func TestRemoveAclAdmins(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := createNAclAdmin(keeper, ctx, 10)

	var addresses []string

	for _, item := range items {
		addresses = append(addresses, item.Address)
	}

	keeper.RemoveAclAdmins(ctx, addresses)

	for i := range addresses {
		_, found := keeper.GetAclAdmin(ctx, addresses[i])
		require.False(t, found)
	}
}

func TestAclAdminGetAll(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := createNAclAdmin(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAclAdmin(ctx)),
	)
}
