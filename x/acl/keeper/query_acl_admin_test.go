package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestAclAdminQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	msgs := createNAclAdmin(keeper, ctx, 2)
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
			response, err := keeper.AclAdmin(ctx, tc.request)
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

func TestAclAdminQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	msgs := createNAclAdmin(keeper, ctx, 5)

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
			resp, err := keeper.AclAdminAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAdmin), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.AclAdmin),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AclAdminAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAdmin), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.AclAdmin),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AclAdminAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.AclAdmin),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AclAdminAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
