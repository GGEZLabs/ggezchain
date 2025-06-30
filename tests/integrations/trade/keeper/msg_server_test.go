package keeper_test

import (
	"testing"
	"time"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gotest.tools/v3/assert"
)

func TestCreateTrade(t *testing.T) {
	f := initFixture(t)

	msgServer := keeper.NewMsgServerImpl(*f.tradeKeeper)

	// Set AclAuthority
	setAclAuthority(f.ctx, f.aclKeeper)

	testCases := []struct {
		name      string
		exceptErr bool
		req       types.MsgCreateTrade
		expErrMsg string
	}{
		{
			name:      "no permission - not found in aclAuthority",
			exceptErr: true,
			req: types.MsgCreateTrade{
				Creator: testutil.Eve,
			},
			expErrMsg: "authority address does not exist",
		},
		{
			name:      "no permission for module",
			exceptErr: true,
			req: types.MsgCreateTrade{
				Creator: testutil.Carol,
			},
			expErrMsg: "no permission for module trade",
		},
		{
			name:      "does not has maker permission",
			exceptErr: true,
			req: types.MsgCreateTrade{
				Creator: testutil.Bob,
			},
			expErrMsg: "invalid maker permission",
		},
		{
			name:      "invalid trade data",
			exceptErr: true,
			req: types.MsgCreateTrade{
				Creator:   testutil.Alice,
				TradeData: `{"trade_info":{"asset_holder_id":0,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			},
			expErrMsg: "invalid trade info",
		},
		{
			name:      "all good",
			exceptErr: false,
			req:       *types.GetSampleMsgCreateTrade(),
			expErrMsg: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := msgServer.CreateTrade(f.ctx, &testCase.req)
			if testCase.exceptErr {
				assert.ErrorContains(t, err, testCase.expErrMsg)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}

func TestCanceledTradeAfterCreateTrade(t *testing.T) {
	f := initFixture(t)
	header := f.ctx.BlockHeader()
	header.Time = time.Now()
	f.ctx = f.ctx.WithBlockHeader(header)

	msgServer := keeper.NewMsgServerImpl(*f.tradeKeeper)

	// Set AclAuthority
	setAclAuthority(f.ctx, f.aclKeeper)

	_, err := msgServer.CreateTrade(f.ctx, types.NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		types.GetSampleTradeData(types.TradeTypeBuy),
		"{}",
		"",
		"",
	),
	)
	assert.NilError(t, err)

	_, err = msgServer.CreateTrade(f.ctx, types.NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		types.GetSampleTradeData(types.TradeTypeBuy),
		"{}",
		"",
		"",
	),
	)
	assert.NilError(t, err)

	_, err = msgServer.CreateTrade(f.ctx, types.NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		types.GetSampleTradeData(types.TradeTypeBuy),
		"{}",
		"",
		"",
	),
	)
	assert.NilError(t, err)

	allTrades := f.tradeKeeper.GetAllStoredTrade(f.ctx)
	assert.Equal(t, len(allTrades), 3)

	allTempTrades := f.tradeKeeper.GetAllStoredTempTrade(f.ctx)
	assert.Equal(t, len(allTempTrades), 3)

	// update create date
	f.tradeKeeper.SetStoredTempTrade(f.ctx, types.StoredTempTrade{
		TradeIndex: allTempTrades[0].TradeIndex,
		TxDate:     "2025-05-02T16:06:05Z",
	})

	f.tradeKeeper.SetStoredTempTrade(f.ctx, types.StoredTempTrade{
		TradeIndex: allTempTrades[1].TradeIndex,
		TxDate:     "2025-05-01T18:09:05Z",
	})

	f.tradeKeeper.SetStoredTempTrade(f.ctx, types.StoredTempTrade{
		TradeIndex: allTempTrades[2].TradeIndex,
		TxDate:     "2025-05-05T20:04:05Z",
	})

	_, err = msgServer.CreateTrade(f.ctx, types.NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		types.GetSampleTradeData(types.TradeTypeBuy),
		"{}",
		"",
		"",
	),
	)
	assert.NilError(t, err)

	allTrades = f.tradeKeeper.GetAllStoredTrade(f.ctx)
	assert.Equal(t, len(allTrades), 4)

	allTempTrades = f.tradeKeeper.GetAllStoredTempTrade(f.ctx)
	assert.Equal(t, len(allTempTrades), 1)

	assert.Equal(t, allTrades[0].Status, types.StatusCanceled)
	assert.Equal(t, allTrades[1].Status, types.StatusCanceled)
	assert.Equal(t, allTrades[2].Status, types.StatusCanceled)
	assert.Equal(t, allTrades[3].Status, types.StatusPending)
}

func TestProcessTrade(t *testing.T) {
	f := initFixture(t)

	msgServer := keeper.NewMsgServerImpl(*f.tradeKeeper)

	// Set AclAuthority
	setAclAuthority(f.ctx, f.aclKeeper)

	// Create trades
	_, err := msgServer.CreateTrade(f.ctx, types.NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		types.GetSampleTradeData(types.TradeTypeBuy),
		"{}",
		"",
		"",
	),
	)
	assert.NilError(t, err)

	_, err = msgServer.CreateTrade(f.ctx, types.NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		types.GetSampleTradeData(types.TradeTypeBuy),
		"{}",
		"",
		"",
	),
	)
	assert.NilError(t, err)
	_, err = msgServer.CreateTrade(f.ctx, types.NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		types.GetSampleTradeData(types.TradeTypeBuy),
		"{}",
		"",
		"",
	),
	)
	assert.NilError(t, err)

	testCases := []struct {
		name      string
		exceptErr bool
		req       types.MsgProcessTrade
		expErrMsg string
	}{
		{
			name:      "no permission - not found in aclAuthority",
			exceptErr: true,
			req: types.MsgProcessTrade{
				Creator: testutil.Eve,
			},
			expErrMsg: "authority address does not exist",
		},
		{
			name:      "no permission for module",
			exceptErr: true,
			req: types.MsgProcessTrade{
				Creator: testutil.Carol,
			},
			expErrMsg: "no permission for module trade",
		},
		{
			name:      "invalid checker permission",
			exceptErr: true,
			req: types.MsgProcessTrade{
				Creator: testutil.Alice,
			},
			expErrMsg: "invalid checker permission",
		},
		{
			name:      "trade index not found",
			exceptErr: true,
			req: types.MsgProcessTrade{
				Creator:    testutil.Bob,
				TradeIndex: 50,
			},
			expErrMsg: "not found",
		},
		{
			name:      "all good",
			exceptErr: false,
			req: types.MsgProcessTrade{
				Creator:     testutil.Bob,
				ProcessType: types.ProcessTypeConfirm,
				TradeIndex:  1,
			},
			expErrMsg: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := msgServer.ProcessTrade(f.ctx, &testCase.req)
			if testCase.exceptErr {
				assert.ErrorContains(t, err, testCase.expErrMsg)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}

func TestSupplyAndBalancesAfterProcessTrade(t *testing.T) {
	f := initFixture(t)

	msgServer := keeper.NewMsgServerImpl(*f.tradeKeeper)

	// Set AclAuthority
	setAclAuthority(f.ctx, f.aclKeeper)

	// Trade 1
	_, err := msgServer.CreateTrade(f.ctx, types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeBuy, 500000000))
	assert.NilError(t, err)

	_, err = msgServer.ProcessTrade(f.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  1,
	})
	assert.NilError(t, err)

	trade, found := f.tradeKeeper.GetStoredTrade(f.ctx, 1)
	assert.Assert(t, found == true)
	assert.Assert(t, trade.Status == types.StatusProcessed)

	supply := f.bankKeeper.GetSupply(f.ctx, types.DefaultDenom)
	assert.Assert(t, supply.Amount.Int64() == 500000000)

	receiverAddress, err := sdk.AccAddressFromBech32(testutil.Alice)
	assert.NilError(t, err)

	balance := f.bankKeeper.GetBalance(f.ctx, receiverAddress, types.DefaultDenom)
	assert.Assert(t, balance.Amount.Int64() == 500000000)

	// Trade 2
	_, err = msgServer.CreateTrade(f.ctx, types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 700000000))

	assert.NilError(t, err)

	_, err = msgServer.ProcessTrade(f.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  2,
	})
	assert.NilError(t, err)

	trade, found = f.tradeKeeper.GetStoredTrade(f.ctx, 2)
	assert.Assert(t, found == true)
	assert.Assert(t, trade.Status == types.StatusRejected)

	supply = f.bankKeeper.GetSupply(f.ctx, types.DefaultDenom)
	assert.Assert(t, supply.Amount.Int64() == 500000000)

	receiverAddress, err = sdk.AccAddressFromBech32(testutil.Alice)
	assert.NilError(t, err)

	balance = f.bankKeeper.GetBalance(f.ctx, receiverAddress, types.DefaultDenom)
	assert.Assert(t, balance.Amount.Int64() == 500000000)

	// Trade 3
	_, err = msgServer.CreateTrade(f.ctx, types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 700000000))
	assert.NilError(t, err)

	_, err = msgServer.ProcessTrade(f.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  3,
	})
	assert.NilError(t, err)

	trade, found = f.tradeKeeper.GetStoredTrade(f.ctx, 3)
	assert.Assert(t, found == true)
	assert.Assert(t, trade.Status == types.StatusFailed)

	supply = f.bankKeeper.GetSupply(f.ctx, types.DefaultDenom)
	assert.Assert(t, supply.Amount.Int64() == 500000000)

	receiverAddress, err = sdk.AccAddressFromBech32(testutil.Alice)
	assert.NilError(t, err)

	balance = f.bankKeeper.GetBalance(f.ctx, receiverAddress, types.DefaultDenom)
	assert.Assert(t, balance.Amount.Int64() == 500000000)

	// Trade 4
	_, err = msgServer.CreateTrade(f.ctx, types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 100000000))
	assert.NilError(t, err)

	_, err = msgServer.ProcessTrade(f.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  4,
	})
	assert.NilError(t, err)

	trade, found = f.tradeKeeper.GetStoredTrade(f.ctx, 4)
	assert.Assert(t, found == true)
	assert.Assert(t, trade.Status == types.StatusProcessed)

	supply = f.bankKeeper.GetSupply(f.ctx, types.DefaultDenom)
	assert.Assert(t, supply.Amount.Int64() == 400000000)

	receiverAddress, err = sdk.AccAddressFromBech32(testutil.Alice)
	assert.NilError(t, err)

	balance = f.bankKeeper.GetBalance(f.ctx, receiverAddress, types.DefaultDenom)
	assert.Assert(t, balance.Amount.Int64() == 400000000)
}
