package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestStoredTradeQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNStoredTrade(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetStoredTradeRequest
		response *types.QueryGetStoredTradeResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetStoredTradeRequest{
				TradeIndex: msgs[0].TradeIndex,
			},
			response: &types.QueryGetStoredTradeResponse{StoredTrade: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetStoredTradeRequest{
				TradeIndex: msgs[1].TradeIndex,
			},
			response: &types.QueryGetStoredTradeResponse{StoredTrade: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetStoredTradeRequest{
				TradeIndex: uint64(100000),
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
			response, err := keeper.StoredTrade(wctx, tc.request)
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

func TestStoredTradeQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNStoredTrade(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllStoredTradeRequest {
		return &types.QueryAllStoredTradeRequest{
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
			resp, err := keeper.StoredTradeAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTrade), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.StoredTrade),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.StoredTradeAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTrade), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.StoredTrade),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.StoredTradeAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.StoredTrade),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.StoredTradeAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
