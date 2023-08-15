package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

func TestTradeIndexQuery(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestTradeIndex(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetTradeIndexRequest
		response *types.QueryGetTradeIndexResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTradeIndexRequest{},
			response: &types.QueryGetTradeIndexResponse{TradeIndex: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.TradeIndex(wctx, tc.request)
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
