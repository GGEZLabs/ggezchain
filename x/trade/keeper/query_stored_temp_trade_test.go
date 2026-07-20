package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createNStoredTempTrade(keeper keeper.Keeper, ctx context.Context, n int) []types.StoredTempTrade {
	items := make([]types.StoredTempTrade, n)
	for i := range items {
		items[i].TradeIndex = uint64(i)
		items[i].TxDate = strconv.Itoa(i)
		_ = keeper.StoredTempTrade.Set(ctx, items[i].TradeIndex, items[i])
	}
	return items
}

func TestStoredTempTradeQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNStoredTempTrade(f.keeper, f.ctx, 2)
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
				TradeIndex: 100000,
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
			response, err := qs.GetStoredTempTrade(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestStoredTempTradeQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNStoredTempTrade(f.keeper, f.ctx, 5)

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
			resp, err := qs.ListStoredTempTrade(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTempTrade), step)
			require.Subset(t, msgs, resp.StoredTempTrade)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListStoredTempTrade(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTempTrade), step)
			require.Subset(t, msgs, resp.StoredTempTrade)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListStoredTempTrade(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.StoredTempTrade)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListStoredTempTrade(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
