package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GGEZLabs/ramichain/x/acl/keeper"
	"github.com/GGEZLabs/ramichain/x/acl/types"
)

func TestSuperAdminQuery(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	item := types.SuperAdmin{}
	err := f.keeper.SuperAdmin.Set(f.ctx, item)
	require.NoError(t, err)

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
			response, err := qs.GetSuperAdmin(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}
