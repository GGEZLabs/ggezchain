package keeper_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestIsSuperAdmin(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	admin := sample.AccAddress()
	alice := sample.AccAddress()

	testCases := []struct {
		name           string
		address        string
		fun            func()
		expectedOutput bool
	}{
		{
			name:           "super admin not found",
			address:        "",
			fun:            func() {},
			expectedOutput: false,
		},
		{
			name:    "address does not match super admin",
			address: alice,
			fun: func() {
				keeper.SetSuperAdmin(ctx, types.SuperAdmin{Admin: admin})
			},
			expectedOutput: false,
		},
		{
			name:           "empty input address",
			address:        "",
			fun:            func() {},
			expectedOutput: false,
		},
		{
			name:           "all good",
			address:        admin,
			fun:            func() {},
			expectedOutput: true,
		},
	}

	for _, tc := range testCases {
		tc.fun()
		t.Run(tc.name, func(t *testing.T) {
			isAdmin := keeper.IsSuperAdmin(ctx, tc.address)
			require.Equal(t, tc.expectedOutput, isAdmin)
		})
	}
}

func TestIsAdmin(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	admin := sample.AccAddress()
	addr := sample.AccAddress()
	keeper.SetAclAdmin(ctx, types.AclAdmin{Address: admin})
	testCases := []struct {
		name           string
		address        string
		expectedOutput bool
	}{
		{
			name:           "address does not match admin",
			address:        addr,
			expectedOutput: false,
		},
		{
			name:           "empty input address",
			address:        "",
			expectedOutput: false,
		},
		{
			name:           "all good",
			address:        admin,
			expectedOutput: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isAdmin := keeper.IsAdmin(ctx, tc.address)
			require.Equal(t, tc.expectedOutput, isAdmin)
		})
	}
}
