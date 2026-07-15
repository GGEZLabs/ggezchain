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
				err := f.tradeKeeper.TradeIndex.Set(ctx, types.TradeIndex{NextId: 1})
				assert.NilError(t, err)
				tradeIndex, err = f.tradeKeeper.TradeIndex.Get(ctx)
				assert.NilError(t, err)
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

			trade, err := queryClient.GetTradeIndex(gocontext.Background(), req)

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
				err := f.tradeKeeper.StoredTrade.Set(ctx, 1, types.StoredTrade{TradeIndex: 1})
				assert.NilError(t, err)
				storedTrade, err = f.tradeKeeper.StoredTrade.Get(ctx, 1)
				assert.NilError(t, err)
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

			trade, err := queryClient.GetStoredTrade(gocontext.Background(), req)

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
				assert.NilError(t, f.tradeKeeper.StoredTrade.Set(ctx, 1, types.StoredTrade{TradeIndex: 1}))
				assert.NilError(t, f.tradeKeeper.StoredTrade.Set(ctx, 2, types.StoredTrade{TradeIndex: 2}))
				assert.NilError(t, f.tradeKeeper.StoredTrade.Set(ctx, 3, types.StoredTrade{TradeIndex: 3}))
				assert.NilError(t, f.tradeKeeper.StoredTrade.Set(ctx, 4, types.StoredTrade{TradeIndex: 4}))
				assert.NilError(t, f.tradeKeeper.StoredTrade.Set(ctx, 5, types.StoredTrade{TradeIndex: 5}))
				storedTradeAll = getAllStoredTrade(t, ctx, *f.tradeKeeper)
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
				storedTradeAll = getAllStoredTrade(t, ctx, *f.tradeKeeper)
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

			trades, err := queryClient.ListStoredTrade(gocontext.Background(), req)

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
				err := f.tradeKeeper.StoredTempTrade.Set(ctx, 1, types.StoredTempTrade{TradeIndex: 1})
				assert.NilError(t, err)
				storedTempTrade, err = f.tradeKeeper.StoredTempTrade.Get(ctx, 1)
				assert.NilError(t, err)
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

			trade, err := queryClient.GetStoredTempTrade(gocontext.Background(), req)

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
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 1, types.StoredTempTrade{TradeIndex: 1}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 2, types.StoredTempTrade{TradeIndex: 2}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 3, types.StoredTempTrade{TradeIndex: 3}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 4, types.StoredTempTrade{TradeIndex: 4}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 5, types.StoredTempTrade{TradeIndex: 5}))
				storedTempTradeAll = getAllStoredTempTrade(t, ctx, *f.tradeKeeper)
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
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 1, types.StoredTempTrade{TradeIndex: 1}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 2, types.StoredTempTrade{TradeIndex: 2}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 3, types.StoredTempTrade{TradeIndex: 3}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 4, types.StoredTempTrade{TradeIndex: 4}))
				assert.NilError(t, f.tradeKeeper.StoredTempTrade.Set(ctx, 5, types.StoredTempTrade{TradeIndex: 5}))
				storedTempTradeAll = getAllStoredTempTrade(t, ctx, *f.tradeKeeper)
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

			trades, err := queryClient.ListStoredTempTrade(gocontext.Background(), req)

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
