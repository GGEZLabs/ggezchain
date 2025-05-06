package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestAclAuthorityQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	msgs := createNAclAuthority(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetAclAuthorityRequest
		response *types.QueryGetAclAuthorityResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAclAuthorityRequest{
				Address: msgs[0].Address,
			},
			response: &types.QueryGetAclAuthorityResponse{AclAuthority: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAclAuthorityRequest{
				Address: msgs[1].Address,
			},
			response: &types.QueryGetAclAuthorityResponse{AclAuthority: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetAclAuthorityRequest{
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
			response, err := keeper.AclAuthority(ctx, tc.request)
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

func TestAclAuthorityQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.AclKeeper(t)
	msgs := createNAclAuthority(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAclAuthorityRequest {
		return &types.QueryAllAclAuthorityRequest{
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
			resp, err := keeper.AclAuthorityAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAuthority), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.AclAuthority),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AclAuthorityAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAuthority), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.AclAuthority),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AclAuthorityAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.AclAuthority),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AclAuthorityAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
