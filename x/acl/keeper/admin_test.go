package keeper_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestIsAdmin(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	admin := sample.AccAddress()
	addr := sample.AccAddress()
	require.NoError(t, keeper.SetParams(ctx, types.Params{Admin: admin}))
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
