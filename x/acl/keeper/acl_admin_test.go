package keeper_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func createNAclAdminItems(k keeper.Keeper, ctx context.Context, n int) []types.AclAdmin {
	items := make([]types.AclAdmin, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)
		_ = k.AclAdmin.Set(ctx, items[i].Address, items[i])
	}
	return items
}

func TestAclAdminSet(t *testing.T) {
	f := initFixture(t)
	items := make([]types.AclAdmin, 10)
	for i := range items {
		items[i].Address = strconv.Itoa(i)
		require.NoError(t, f.keeper.AclAdmin.Set(f.ctx, items[i].Address, items[i]))
	}
	for _, item := range items {
		rst, err := f.keeper.AclAdmin.Get(f.ctx, item.Address)
		require.NoError(t, err)
		require.EqualExportedValues(t, item, rst)
	}
}

func TestAclAdminGet(t *testing.T) {
	f := initFixture(t)
	items := createNAclAdminItems(f.keeper, f.ctx, 10)
	for _, item := range items {
		rst, err := f.keeper.AclAdmin.Get(f.ctx, item.Address)
		require.NoError(t, err)
		require.EqualExportedValues(t, item, rst)
	}
}

func TestAclAdminRemove(t *testing.T) {
	f := initFixture(t)
	items := createNAclAdminItems(f.keeper, f.ctx, 10)
	for _, item := range items {
		require.NoError(t, f.keeper.AclAdmin.Remove(f.ctx, item.Address))
		_, err := f.keeper.AclAdmin.Get(f.ctx, item.Address)
		require.Error(t, err)
		require.True(t, errors.Is(err, collections.ErrNotFound))
	}
}

func TestRemoveAclAdmins(t *testing.T) {
	f := initFixture(t)
	items := createNAclAdminItems(f.keeper, f.ctx, 10)

	var addresses []string
	for _, item := range items {
		addresses = append(addresses, item.Address)
	}

	for _, address := range addresses {
		require.NoError(t, f.keeper.AclAdmin.Remove(f.ctx, address))
	}

	for i := range addresses {
		_, err := f.keeper.AclAdmin.Get(f.ctx, addresses[i])
		require.Error(t, err)
		require.True(t, errors.Is(err, collections.ErrNotFound))
	}
}

func TestAclAdminGetAll(t *testing.T) {
	f := initFixture(t)
	items := createNAclAdminItems(f.keeper, f.ctx, 10)
	all, err := f.keeper.GetAllAclAdmin(f.ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, items, all)
}
