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

func TestStoredTempTradeQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNStoredTempTrade(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetStoredTempTradeRequest
		response *types.QueryGetStoredTempTradeResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetStoredTempTradeRequest{
				TradeIndex: msgs[0].TradeIndex,
			},
			response: &types.QueryGetStoredTempTradeResponse{StoredTempTrade: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetStoredTempTradeRequest{
				TradeIndex: msgs[1].TradeIndex,
			},
			response: &types.QueryGetStoredTempTradeResponse{StoredTempTrade: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetStoredTempTradeRequest{
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
			response, err := keeper.StoredTempTrade(wctx, tc.request)
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

func TestStoredTempTradeQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNStoredTempTrade(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllStoredTempTradeRequest {
		return &types.QueryAllStoredTempTradeRequest{
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
			resp, err := keeper.StoredTempTradeAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTempTrade), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.StoredTempTrade),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.StoredTempTradeAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTempTrade), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.StoredTempTrade),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.StoredTempTradeAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.StoredTempTrade),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.StoredTempTradeAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
