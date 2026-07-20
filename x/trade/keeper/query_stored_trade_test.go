package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createNStoredTrade(keeper keeper.Keeper, ctx context.Context, n int) []types.StoredTrade {
	items := make([]types.StoredTrade, n)
	for i := range items {
		items[i].TradeIndex = uint64(i)
		items[i].TradeType = types.TradeType(i)
		items[i].Amount = sdk.NewInt64Coin(`token`, int64(i+100))
		items[i].CoinMintingPriceUsd = strconv.Itoa(i)
		items[i].ReceiverAddress = strconv.Itoa(i)
		items[i].Status = types.TradeStatus(i)
		items[i].Maker = strconv.Itoa(i)
		items[i].Checker = strconv.Itoa(i)
		items[i].TxDate = strconv.Itoa(i)
		items[i].CreateDate = strconv.Itoa(i)
		items[i].UpdateDate = strconv.Itoa(i)
		items[i].ProcessDate = strconv.Itoa(i)
		items[i].TradeData = strconv.Itoa(i)
		items[i].CoinMintingPriceJson = strconv.Itoa(i)
		items[i].ExchangeRateJson = strconv.Itoa(i)
		items[i].BankingSystemData = strconv.Itoa(i)
		items[i].Result = strconv.Itoa(i)
		_ = keeper.StoredTrade.Set(ctx, items[i].TradeIndex, items[i])
	}
	return items
}

func TestStoredTradeQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNStoredTrade(f.keeper, f.ctx, 2)
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
			response, err := qs.GetStoredTrade(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestStoredTradeQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNStoredTrade(f.keeper, f.ctx, 5)

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
			resp, err := qs.ListStoredTrade(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTrade), step)
			require.Subset(t, msgs, resp.StoredTrade)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListStoredTrade(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredTrade), step)
			require.Subset(t, msgs, resp.StoredTrade)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListStoredTrade(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.StoredTrade)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListStoredTrade(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
