package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestProcessTradeConfirm() {
	indexes := suite.createTrade(1)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)
	suite.EqualValues(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusProcessed,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestProcessTradeReject() {
	indexes := suite.createTrade(1)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)
	suite.EqualValues(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusRejected,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestProcessTradeWithInvalidCheckerPermission() {
	indexes := suite.createTrade(1)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Alice,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.ErrorIs(err, types.ErrInvalidCheckerPermission)
}

func (suite *KeeperTestSuite) TestStoredTradeAfterConfirmTrade() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTradeConfirmed(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestTempTradeAfterConfirmTrade() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     1,
		CreateDate:     types.GetSampleStoredTrade(indexes[0]).CreateDate,
		TempTradeIndex: 1,
	}, tempTrade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	// should remove temp trade after process
	tempTrade, found = keeper.GetStoredTempTrade(suite.ctx, indexes[0])
	suite.False(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     0,
		CreateDate:     "",
		TempTradeIndex: 0,
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestStoredTradeAfterRejectTrade() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTradeRejected(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestTempTradeAfterRejectTrade() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, 1)

	suite.False(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     0,
		CreateDate:     "",
		TempTradeIndex: 0,
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestProcessTwoTrade() {
	indexes := suite.createTrade(2)
	keeper := suite.app.TradeKeeper

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrades := keeper.GetAllStoredTempTrade(suite.ctx)
	trades := keeper.GetAllStoredTrade(suite.ctx)
	suite.EqualValues(len(tempTrades), len(indexes))
	suite.EqualValues(len(trades), len(indexes))

	for _, tradeIndex := range indexes {

		_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
			Creator:     testutil.Bob,
			ProcessType: types.ProcessTypeConfirm,
			TradeIndex:  tradeIndex,
		})
		suite.Nil(err)
	}

	tempTrades = keeper.GetAllStoredTempTrade(suite.ctx)
	trades = keeper.GetAllStoredTrade(suite.ctx)
	suite.EqualValues(len(tempTrades), 0)
	suite.EqualValues(len(trades), len(indexes))
}

func (suite *KeeperTestSuite) TestProcessTradeInsufficientFund() {
	suite.setupTest()
	keeper := suite.app.TradeKeeper

	msgCreateTradeBuy := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeBuy, 5000000000000)
	msgCreateTradeSell := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 7000000000000)

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeBuy)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	createResponse, err = suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeSell)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	trade, found := keeper.GetStoredTrade(suite.ctx, 2)
	suite.True(found)
	suite.EqualValues(types.StatusFailed, trade.Status)
}

func (suite *KeeperTestSuite) TestProcessTradeAlreadyConfirmed() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTradeConfirmed(indexes[0]), trade)

	// process confirmed trade
	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.ErrorIs(err, types.ErrTradeStatusCompleted)
}

func (suite *KeeperTestSuite) TestProcessTradeAlreadyRejected() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.True(found)
	suite.EqualValues(types.GetSampleStoredTradeRejected(indexes[0]), trade)

	// process confirmed trade
	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.ErrorIs(err, types.ErrTradeStatusRejected)
}

func (suite *KeeperTestSuite) TestProcessTradeCheckerIsNotMaker() {
	suite.setupTest()
	msg := types.GetSampleMsgCreateTrade()
	msg.Creator = testutil.Trent
	_, err := suite.msgServer.CreateTrade(suite.ctx, msg)
	suite.Nil(err)

	keeper := suite.app.TradeKeeper

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Trent,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  1,
	})
	suite.ErrorIs(err, types.ErrCheckerMustBeDifferent)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)

	trade, found := keeper.GetStoredTrade(suite.ctx, 1)
	suite.True(found)
	suite.EqualValues(types.StatusPending, trade.Status)
}

func (suite *KeeperTestSuite) TestSupplyAfterProcessTrade() {
	suite.setupTest()
	msgCreateTradeBuy := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeBuy, 5000000000000)
	msgCreateTradeSell := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 3000000000000)

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeBuy)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	supply := suite.app.BankKeeper.GetSupply(suite.ctx, types.DefaultCoinDenom)
	suite.Require().Equal(sdkmath.NewInt(5000000000000), supply.Amount)

	createResponse, err = suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeSell)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	supply = suite.app.BankKeeper.GetSupply(suite.ctx, types.DefaultCoinDenom)
	suite.Require().Equal(sdkmath.NewInt(2000000000000), supply.Amount)

}

func (suite *KeeperTestSuite) TestBalancesAfterProcessTrade() {

	suite.setupTest()
	msgCreateTradeBuy := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeBuy, 5000000000000)
	msgCreateTradeSell := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 3000000000000)

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeBuy)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	receiverAddress, err := sdk.AccAddressFromBech32(msgCreateTradeBuy.ReceiverAddress)
	suite.Require().NoError(err)

	initialBalance := suite.app.BankKeeper.GetBalance(suite.ctx, receiverAddress, types.DefaultCoinDenom)
	suite.Require().Equal(sdkmath.NewInt(5000000000000), initialBalance.Amount)

	createResponse, err = suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeSell)
	suite.Require().NoError(err)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	finalBalance := suite.app.BankKeeper.GetBalance(suite.ctx, receiverAddress, types.DefaultCoinDenom)

	suite.Require().Equal(sdkmath.NewInt(2000000000000), finalBalance.Amount)

}
