package keeper_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSuperAdminQuery(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	item := createTestSuperAdmin(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetSuperAdminRequest
		response *types.QueryGetSuperAdminResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetSuperAdminRequest{},
			response: &types.QueryGetSuperAdminResponse{SuperAdmin: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.SuperAdmin(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
