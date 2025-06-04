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

func createNAclAuthority(keeper keeper.Keeper, ctx context.Context, n int) []types.AclAuthority {
	items := make([]types.AclAuthority, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)
		items[i].Name = strconv.Itoa(i)
		items[i].AccessDefinitions = []*types.AccessDefinition{
			{Module: strconv.Itoa(i), IsMaker: true, IsChecker: false},
		}
		_ = keeper.AclAuthority.Set(ctx, items[i].Address, items[i])
	}
	return items
}

func TestAclAuthorityQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAclAuthority(f.keeper, f.ctx, 2)
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
			response, err := qs.GetAclAuthority(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestAclAuthorityQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAclAuthority(f.keeper, f.ctx, 5)

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
			resp, err := qs.ListAclAuthority(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAuthority), step)
			require.Subset(t, msgs, resp.AclAuthority)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAclAuthority(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AclAuthority), step)
			require.Subset(t, msgs, resp.AclAuthority)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListAclAuthority(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.AclAuthority)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListAclAuthority(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
