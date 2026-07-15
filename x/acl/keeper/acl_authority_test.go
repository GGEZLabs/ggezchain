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

func createNAclAuthorityItems(k keeper.Keeper, ctx context.Context, n int) []types.AclAuthority {
	items := make([]types.AclAuthority, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)
		_ = k.AclAuthority.Set(ctx, items[i].Address, items[i])
	}
	return items
}

func TestAclAuthorityGet(t *testing.T) {
	f := initFixture(t)
	items := createNAclAuthorityItems(f.keeper, f.ctx, 10)
	for _, item := range items {
		rst, err := f.keeper.AclAuthority.Get(f.ctx, item.Address)
		require.NoError(t, err)
		require.EqualExportedValues(t, item, rst)
	}
}

func TestAclAuthorityRemove(t *testing.T) {
	f := initFixture(t)
	items := createNAclAuthorityItems(f.keeper, f.ctx, 10)
	for _, item := range items {
		require.NoError(t, f.keeper.AclAuthority.Remove(f.ctx, item.Address))
		_, err := f.keeper.AclAuthority.Get(f.ctx, item.Address)
		require.Error(t, err)
		require.True(t, errors.Is(err, collections.ErrNotFound))
	}
}

func TestAclAuthorityGetAll(t *testing.T) {
	f := initFixture(t)
	items := createNAclAuthorityItems(f.keeper, f.ctx, 10)

	var all []types.AclAuthority
	err := f.keeper.AclAuthority.Walk(f.ctx, nil, func(_ string, val types.AclAuthority) (stop bool, err error) {
		all = append(all, val)
		return false, nil
	})
	require.NoError(t, err)
	require.ElementsMatch(t, items, all)
}
