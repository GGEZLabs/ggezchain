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

func createNAclAuthority(keeper keeper.Keeper, ctx context.Context, n int) []types.AclAuthority {
	items := make([]types.AclAuthority, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)

		keeper.SetAclAuthority(ctx, items[i])
	}
	return items
}

func TestAclAuthorityGet(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := createNAclAuthority(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetAclAuthority(ctx,
			item.Address,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestAclAuthorityRemove(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := createNAclAuthority(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAclAuthority(ctx,
			item.Address,
		)
		_, found := keeper.GetAclAuthority(ctx,
			item.Address,
		)
		require.False(t, found)
	}
}

func TestAclAuthorityGetAll(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	items := createNAclAuthority(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAclAuthority(ctx)),
	)
}
