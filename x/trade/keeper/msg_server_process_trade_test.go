package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gomock "go.uber.org/mock/gomock"
)

func (suite *KeeperTestSuite) TestProcessTradeConfirm() {
	indexes := suite.createNTrades(1)

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)
	suite.Require().EqualValues(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusProcessed,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestProcessTradeReject() {
	indexes := suite.createNTrades(1)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)
	suite.Require().EqualValues(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusRejected,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestProcessTradeWithInvalidCheckerPermission() {
	indexes := suite.createNTrades(1)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Alice,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().ErrorIs(err, types.ErrInvalidCheckerPermission)
}

func (suite *KeeperTestSuite) TestStoredTradeAfterConfirmTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTradeConfirmed(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestTempTradeAfterConfirmTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTempTrade{
		TradeIndex:     1,
		CreateDate:     types.GetSampleStoredTrade(indexes[0]).CreateDate,
		TempTradeIndex: 1,
	}, tempTrade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	// should remove temp trade after process
	tempTrade, found = keeper.GetStoredTempTrade(suite.ctx, indexes[0])
	suite.Require().False(found)
	suite.Require().EqualValues(types.StoredTempTrade{
		TradeIndex:     0,
		CreateDate:     "",
		TempTradeIndex: 0,
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestStoredTradeAfterRejectTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTradeRejected(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestTempTradeAfterRejectTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, 1)

	suite.Require().False(found)
	suite.Require().EqualValues(types.StoredTempTrade{
		TradeIndex:     0,
		CreateDate:     "",
		TempTradeIndex: 0,
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestProcessTwoTrade() {
	indexes := suite.createNTrades(2)
	keeper := suite.tradeKeeper

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(2)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(2)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrades := keeper.GetAllStoredTempTrade(suite.ctx)
	trades := keeper.GetAllStoredTrade(suite.ctx)
	suite.Require().EqualValues(len(tempTrades), len(indexes))
	suite.Require().EqualValues(len(trades), len(indexes))

	for _, tradeIndex := range indexes {
		_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
			Creator:     testutil.Bob,
			ProcessType: types.ProcessTypeConfirm,
			TradeIndex:  tradeIndex,
		})
		suite.Require().Nil(err)
	}

	tempTrades = keeper.GetAllStoredTempTrade(suite.ctx)
	trades = keeper.GetAllStoredTrade(suite.ctx)
	suite.Require().EqualValues(len(tempTrades), 0)
	suite.Require().EqualValues(len(trades), len(indexes))
}

func (suite *KeeperTestSuite) TestProcessTradeInsufficientFund() {
	suite.setupTest()
	keeper := suite.tradeKeeper

	msgCreateTradeBuy := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeBuy, 5000000000000)
	msgCreateTradeSell := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 7000000000000)

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeBuy)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	createResponse, err = suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeSell)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(sdkerrors.ErrInsufficientFunds).Times(1)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	trade, found := keeper.GetStoredTrade(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().EqualValues(types.StatusFailed, trade.Status)
}

func (suite *KeeperTestSuite) TestProcessTradeAlreadyConfirmed() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTradeConfirmed(indexes[0]), trade)

	// process confirmed trade
	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().ErrorIs(err, types.ErrInvalidTradeStatus)
}

func (suite *KeeperTestSuite) TestProcessTradeAlreadyRejected() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().Nil(err)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTradeRejected(indexes[0]), trade)

	// process confirmed trade
	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().ErrorIs(err, types.ErrInvalidTradeStatus)
}

func (suite *KeeperTestSuite) TestProcessTradeCheckerIsNotMaker() {
	suite.setupTest()

	msg := types.GetSampleMsgCreateTrade()
	msg.Creator = testutil.Trent
	_, err := suite.msgServer.CreateTrade(suite.ctx, msg)
	suite.Require().Nil(err)

	keeper := suite.tradeKeeper

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Trent,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  1,
	})
	suite.Require().ErrorIs(err, types.ErrCheckerMustBeDifferent)

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)

	trade, found := keeper.GetStoredTrade(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().EqualValues(types.StatusPending, trade.Status)
}

func (suite *KeeperTestSuite) TestSupplyAfterProcessTrade() {
	suite.setupTest()

	msgCreateTradeBuy := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeBuy, 5000000000000)
	msgCreateTradeSell := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 3000000000000)

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeBuy)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().GetSupply(suite.ctx, types.DefaultDenom).Return(sdk.Coin{
		Denom:  types.DefaultDenom,
		Amount: sdkmath.NewInt(5000000000000),
	}).Times(1)

	supply := suite.bankKeeper.GetSupply(suite.ctx, types.DefaultDenom)
	suite.Require().Equal(sdkmath.NewInt(5000000000000), supply.Amount)

	createResponse, err = suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeSell)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().BurnCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().GetSupply(suite.ctx, types.DefaultDenom).Return(sdk.Coin{
		Denom:  types.DefaultDenom,
		Amount: sdkmath.NewInt(2000000000000),
	}).Times(1)

	supply = suite.bankKeeper.GetSupply(suite.ctx, types.DefaultDenom)
	suite.Require().Equal(sdkmath.NewInt(2000000000000), supply.Amount)
}

func (suite *KeeperTestSuite) TestBalancesAfterProcessTrade() {
	suite.setupTest()

	msgCreateTradeBuy := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeBuy, 5000000000000)
	msgCreateTradeSell := types.GetMsgCreateTradeWithTypeAndAmount(types.TradeTypeSell, 3000000000000)

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeBuy)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	receiverAddress, err := sdk.AccAddressFromBech32(msgCreateTradeBuy.ReceiverAddress)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().GetBalance(suite.ctx, gomock.Any(), types.DefaultDenom).Return(sdk.Coin{
		Denom:  types.DefaultDenom,
		Amount: sdkmath.NewInt(5000000000000),
	}).Times(1)

	initialBalance := suite.bankKeeper.GetBalance(suite.ctx, receiverAddress, types.DefaultDenom)
	suite.Require().Equal(sdkmath.NewInt(5000000000000), initialBalance.Amount)

	createResponse, err = suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeSell)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().BurnCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     testutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().GetBalance(suite.ctx, gomock.Any(), types.DefaultDenom).Return(sdk.Coin{
		Denom:  types.DefaultDenom,
		Amount: sdkmath.NewInt(2000000000000),
	}).Times(1)

	finalBalance := suite.bankKeeper.GetBalance(suite.ctx, receiverAddress, types.DefaultDenom)

	suite.Require().Equal(sdkmath.NewInt(2000000000000), finalBalance.Amount)
}
