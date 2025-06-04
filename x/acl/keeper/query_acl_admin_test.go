package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GGEZLabs/ramichain/x/acl/keeper"
	"github.com/GGEZLabs/ramichain/x/acl/types"
)

func createNAclAdmin(keeper keeper.Keeper, ctx context.Context, n int) []types.AclAdmin {
	items := make([]types.AclAdmin, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)
		_ = keeper.AclAdmin.Set(ctx, items[i].Address, items[i])
	}
	return items
}

func TestAclAdminQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAclAdmin(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetAclAdminRequest
		response *types.QueryGetAclAdminResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAclAdminRequest{
				Address: msgs[0].Address,
			},
			response: &types.QueryGetAclAdminResponse{AclAdmin: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAclAdminRequest{
				Address: msgs[1].Address,
			},
			response: &types.QueryGetAclAdminResponse{AclAdmin: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetAclAdminRequest{
				Address: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetAclAdmin(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestAclAdminQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAclAdmin(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAclAdminRequest {
		return &types.QueryAllAclAdminRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAclAdmin(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAdmin), step)
			require.Subset(t, msgs, resp.AclAdmin)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAclAdmin(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAdmin), step)
			require.Subset(t, msgs, resp.AclAdmin)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListAclAdmin(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.AclAdmin)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListAclAdmin(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
