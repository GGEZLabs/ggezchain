package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"gotest.tools/v3/assert"
)

func TestGRPCQueryTradeIndex(t *testing.T) {
	f := initFixture(t)

	ctx, queryClient := f.ctx, f.queryClient

	var (
		req        *types.QueryGetTradeIndexRequest
		res        *types.QueryGetTradeIndexResponse
		tradeIndex types.TradeIndex
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"get trade index",
			func() {
				f.tradeKeeper.SetTradeIndex(ctx, types.TradeIndex{NextId: 1})
				var found bool
				tradeIndex, found = f.tradeKeeper.GetTradeIndex(ctx)
				assert.Assert(t, found == true)
				assert.Assert(t, tradeIndex.String() != "")

				req = &types.QueryGetTradeIndexRequest{}

				res = &types.QueryGetTradeIndexResponse{
					TradeIndex: types.TradeIndex{
						NextId: 1,
					},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			trade, err := queryClient.TradeIndex(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, res.String(), trade.String())
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, trade == nil)
			}
		})
	}
}

func TestGRPCQueryStoredTrade(t *testing.T) {
	f := initFixture(t)

	ctx, queryClient := f.ctx, f.queryClient

	var (
		req         *types.QueryGetStoredTradeRequest
		res         *types.QueryGetStoredTradeResponse
		storedTrade types.StoredTrade
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"empty request",
			func() {
				req = &types.QueryGetStoredTradeRequest{}
			},
			false,
			"not found",
		},
		{
			"zero trade index request",
			func() {
				req = &types.QueryGetStoredTradeRequest{TradeIndex: 0}
			},
			false,
			"not found",
		},
		{
			"query non existed trade",
			func() {
				req = &types.QueryGetStoredTradeRequest{TradeIndex: 1}
			},
			false,
			"not found",
		},
		{
			"store and get trade",
			func() {
				f.tradeKeeper.SetStoredTrade(ctx, types.StoredTrade{TradeIndex: 1})
				var found bool
				storedTrade, found = f.tradeKeeper.GetStoredTrade(ctx, 1)
				assert.Assert(t, found == true)
				assert.Assert(t, storedTrade.String() != "")

				req = &types.QueryGetStoredTradeRequest{TradeIndex: storedTrade.TradeIndex}

				res = &types.QueryGetStoredTradeResponse{
					StoredTrade: types.StoredTrade{TradeIndex: 1},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			trade, err := queryClient.StoredTrade(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, res.String(), trade.String())
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, trade == nil)
			}
		})
	}
}

func TestGRPCQueryAllStoredTrade(t *testing.T) {
	f := initFixture(t)

	ctx, queryClient := f.ctx, f.queryClient

	var (
		req            *types.QueryAllStoredTradeRequest
		res            *types.QueryAllStoredTradeResponse
		storedTradeAll []types.StoredTrade
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"nil request",
			func() {
				req = nil
				res = &types.QueryAllStoredTradeResponse{
					StoredTrade: []types.StoredTrade{},
					Pagination:  &query.PageResponse{},
				}
			},
			true,
			"",
		},
		{
			"store and get all trades",
			func() {
				f.tradeKeeper.SetStoredTrade(ctx, types.StoredTrade{TradeIndex: 1})
				f.tradeKeeper.SetStoredTrade(ctx, types.StoredTrade{TradeIndex: 2})
				f.tradeKeeper.SetStoredTrade(ctx, types.StoredTrade{TradeIndex: 3})
				f.tradeKeeper.SetStoredTrade(ctx, types.StoredTrade{TradeIndex: 4})
				f.tradeKeeper.SetStoredTrade(ctx, types.StoredTrade{TradeIndex: 5})
				storedTradeAll = f.tradeKeeper.GetAllStoredTrade(ctx)
				assert.Assert(t, len(storedTradeAll) == 5)

				req = &types.QueryAllStoredTradeRequest{}

				res = &types.QueryAllStoredTradeResponse{
					StoredTrade: []types.StoredTrade{
						{
							TradeIndex: 1,
						},
						{
							TradeIndex: 2,
						},
						{
							TradeIndex: 3,
						},
						{
							TradeIndex: 4,
						},
						{
							TradeIndex: 5,
						},
					},
					Pagination: &query.PageResponse{
						Total: 5,
					},
				}
			},
			true,
			"",
		},
		{
			"get some of trades",
			func() {
				storedTradeAll = f.tradeKeeper.GetAllStoredTrade(ctx)
				assert.Assert(t, len(storedTradeAll) == 5)

				req = &types.QueryAllStoredTradeRequest{
					Pagination: &query.PageRequest{
						Limit: 2,
					},
				}

				res = &types.QueryAllStoredTradeResponse{
					StoredTrade: []types.StoredTrade{
						{
							TradeIndex: 1,
						},
						{
							TradeIndex: 2,
						},
					},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			trades, err := queryClient.StoredTradeAll(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, len(res.StoredTrade), len(trades.StoredTrade))
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, trades == nil)
			}
		})
	}
}

func TestGRPCQueryStoredTempTrade(t *testing.T) {
	f := initFixture(t)

	ctx, queryClient := f.ctx, f.queryClient

	var (
		req             *types.QueryGetStoredTempTradeRequest
		res             *types.QueryGetStoredTempTradeResponse
		storedTempTrade types.StoredTempTrade
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"empty request",
			func() {
				req = &types.QueryGetStoredTempTradeRequest{}
			},
			false,
			"not found",
		},
		{
			"zero trade index request",
			func() {
				req = &types.QueryGetStoredTempTradeRequest{TradeIndex: 0}
			},
			false,
			"not found",
		},
		{
			"query non existed trade",
			func() {
				req = &types.QueryGetStoredTempTradeRequest{TradeIndex: 1}
			},
			false,
			"not found",
		},
		{
			"store and get trade",
			func() {
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 1})
				var found bool
				storedTempTrade, found = f.tradeKeeper.GetStoredTempTrade(ctx, 1)
				assert.Assert(t, found == true)
				assert.Assert(t, storedTempTrade.String() != "")

				req = &types.QueryGetStoredTempTradeRequest{TradeIndex: storedTempTrade.TradeIndex}

				res = &types.QueryGetStoredTempTradeResponse{
					StoredTempTrade: types.StoredTempTrade{TradeIndex: 1},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			trade, err := queryClient.StoredTempTrade(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, res.String(), trade.String())
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, trade == nil)
			}
		})
	}
}

func TestGRPCQueryAllStoredTempTrade(t *testing.T) {
	f := initFixture(t)

	ctx, queryClient := f.ctx, f.queryClient

	var (
		req                *types.QueryAllStoredTempTradeRequest
		res                *types.QueryAllStoredTempTradeResponse
		storedTempTradeAll []types.StoredTempTrade
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"nil request",
			func() {
				req = nil
				res = &types.QueryAllStoredTempTradeResponse{
					StoredTempTrade: []types.StoredTempTrade{},
					Pagination:      &query.PageResponse{},
				}
			},
			true,
			"",
		},
		{
			"store and get all trades",
			func() {
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 1})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 2})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 3})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 4})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 5})
				storedTempTradeAll = f.tradeKeeper.GetAllStoredTempTrade(ctx)
				assert.Assert(t, len(storedTempTradeAll) == 5)

				req = &types.QueryAllStoredTempTradeRequest{}

				res = &types.QueryAllStoredTempTradeResponse{
					StoredTempTrade: []types.StoredTempTrade{
						{
							TradeIndex: 1,
						},
						{
							TradeIndex: 2,
						},
						{
							TradeIndex: 3,
						},
						{
							TradeIndex: 4,
						},
						{
							TradeIndex: 5,
						},
					},
					Pagination: &query.PageResponse{
						Total: 5,
					},
				}
			},
			true,
			"",
		},
		{
			"get some of trades",
			func() {
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 1})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 2})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 3})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 4})
				f.tradeKeeper.SetStoredTempTrade(ctx, types.StoredTempTrade{TradeIndex: 5})
				storedTempTradeAll = f.tradeKeeper.GetAllStoredTempTrade(ctx)
				assert.Assert(t, len(storedTempTradeAll) == 5)

				req = &types.QueryAllStoredTempTradeRequest{
					Pagination: &query.PageRequest{
						Limit: 2,
					},
				}

				res = &types.QueryAllStoredTempTradeResponse{
					StoredTempTrade: []types.StoredTempTrade{
						{
							TradeIndex: 1,
						},
						{
							TradeIndex: 2,
						},
					},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			trades, err := queryClient.StoredTempTradeAll(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, len(res.StoredTempTrade), len(trades.StoredTempTrade))
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, trades == nil)
			}
		})
	}
}
