package keeper_test

import (
	tradetestutil "github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gomock "go.uber.org/mock/gomock"
)

func (suite *KeeperTestSuite) TestProcessTradeConfirm() {
	indexes := suite.createNTrades(1)

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusProcessed,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestProcessTradeReject() {
	indexes := suite.createNTrades(1)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusRejected,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestProcessTradeWithInvalidCheckerPermission() {
	indexes := suite.createNTrades(1)

	_, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Alice,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().ErrorIs(err, types.ErrInvalidCheckerPermission)
}

// TestProcessTradeValidation covers the former message_process_trade_test.go
// ValidateBasic() cases. That validation now lives inline in the keeper's
// ProcessTrade handler (there's no more standalone MsgProcessTrade.ValidateBasic
// in the new scaffold), so these are ported as keeper-level cases instead.
func (suite *KeeperTestSuite) TestStoredTradeAfterConfirmTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	trade, found := getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTrade(indexes[0]), trade)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusProcessed,
	}, *processResponse)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTradeConfirmed(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestTempTradeAfterConfirmTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	tempTrade, found := getStoredTempTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.StoredTempTrade{
		TradeIndex: 1,
		TxDate:     types.GetSampleStoredTrade(indexes[0]).TxDate,
	}, tempTrade)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusProcessed,
	}, *processResponse)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)

	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	// should remove temp trade after process
	tempTrade, found = getStoredTempTrade(suite.ctx, keeper, indexes[0])
	suite.Require().False(found)
	suite.Require().Equal(types.StoredTempTrade{
		TradeIndex: 0,
		TxDate:     "",
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestStoredTradeAfterRejectTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	trade, found := getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTrade(indexes[0]), trade)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusRejected,
	}, *processResponse)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)

	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTradeRejected(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestTempTradeAfterRejectTrade() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusRejected,
	}, *processResponse)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)

	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrade, found := getStoredTempTrade(suite.ctx, keeper, 1)

	suite.Require().False(found)
	suite.Require().Equal(types.StoredTempTrade{
		TradeIndex: 0,
		TxDate:     "",
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestProcessTwoTrade() {
	indexes := suite.createNTrades(2)
	keeper := suite.tradeKeeper

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(2)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(2)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)

	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrades := getAllStoredTempTrade(suite.ctx, keeper)
	trades := getAllStoredTrade(suite.ctx, keeper)
	suite.Require().Len(indexes, len(tempTrades))
	suite.Require().Len(indexes, len(trades))

	for _, tradeIndex := range indexes {
		processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
			Creator:     tradetestutil.Bob,
			ProcessType: types.ProcessTypeConfirm,
			TradeIndex:  tradeIndex,
		})
		suite.Require().NoError(err)
		suite.Require().Equal(types.MsgProcessTradeResponse{
			TradeIndex: tradeIndex,
			Status:     types.StatusProcessed,
		}, *processResponse)
	}

	tempTrades = getAllStoredTempTrade(suite.ctx, keeper)
	trades = getAllStoredTrade(suite.ctx, keeper)
	suite.Require().Empty(tempTrades)
	suite.Require().Len(indexes, len(trades))
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

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: createResponse.TradeIndex,
		Status:     types.StatusProcessed,
	}, *processResponse)

	createResponse, err = suite.msgServer.CreateTrade(suite.ctx, msgCreateTradeSell)
	suite.Require().NoError(err)

	suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(sdkerrors.ErrInsufficientFunds).Times(1)

	processResponse, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: createResponse.TradeIndex,
		Status:     types.StatusFailed,
	}, *processResponse)

	trade, found := getStoredTrade(suite.ctx, keeper, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.StatusFailed, trade.Status)
}

func (suite *KeeperTestSuite) TestProcessTradeAlreadyConfirmed() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	trade, found := getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTrade(indexes[0]), trade)

	suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName, gomock.Any()).Return(nil).Times(1)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, gomock.Any(), gomock.Any()).Return(nil).Times(1)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusProcessed,
	}, *processResponse)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTradeConfirmed(indexes[0]), trade)

	// process confirmed trade
	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  indexes[0],
	})
	suite.Require().ErrorIs(err, types.ErrInvalidTradeStatus)
}

func (suite *KeeperTestSuite) TestProcessTradeAlreadyRejected() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	trade, found := getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTrade(indexes[0]), trade)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: indexes[0],
		Status:     types.StatusRejected,
	}, *processResponse)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found = getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().Equal(types.GetSampleStoredTradeRejected(indexes[0]), trade)

	// process rejected trade
	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  indexes[0],
	})
	suite.Require().ErrorIs(err, types.ErrInvalidTradeStatus)
}

func (suite *KeeperTestSuite) TestProcessTradeCheckerIsNotMaker() {
	suite.setupTest()

	msg := types.GetSampleMsgCreateTrade()
	msg.Creator = tradetestutil.Trent
	_, err := suite.msgServer.CreateTrade(suite.ctx, msg)
	suite.Require().NoError(err)

	keeper := suite.tradeKeeper

	_, err = suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Trent,
		ProcessType: types.ProcessTypeReject,
		TradeIndex:  1,
	})
	suite.Require().ErrorIs(err, types.ErrCheckerMustBeDifferent)

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().Equal(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)

	trade, found := getStoredTrade(suite.ctx, keeper, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.StatusPending, trade.Status)
}

func (suite *KeeperTestSuite) TestSupplyAfterProcessTradeWithTypeSplit() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeSplit),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})
	suite.Require().NoError(err)

	// No bank keeper calls expected: split trades don't mint/burn.
	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: createResponse.TradeIndex,
		Status:     types.StatusProcessed,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestSupplyAfterProcessTradeWithTypeReverseSplit() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeReverseSplit),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})
	suite.Require().NoError(err)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: createResponse.TradeIndex,
		Status:     types.StatusProcessed,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestSupplyAfterProcessTradeWithTypeReinvestment() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeReinvestment),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})
	suite.Require().NoError(err)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: createResponse.TradeIndex,
		Status:     types.StatusProcessed,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestSupplyAfterProcessTradeWithTypeDividends() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeDividends),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})
	suite.Require().NoError(err)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: createResponse.TradeIndex,
		Status:     types.StatusProcessed,
	}, *processResponse)
}

func (suite *KeeperTestSuite) TestSupplyAfterProcessTradeWithTypeDividendsDeduction() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeDividendsDeduction),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})
	suite.Require().NoError(err)

	processResponse, err := suite.msgServer.ProcessTrade(suite.ctx, &types.MsgProcessTrade{
		Creator:     tradetestutil.Bob,
		ProcessType: types.ProcessTypeConfirm,
		TradeIndex:  createResponse.TradeIndex,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(types.MsgProcessTradeResponse{
		TradeIndex: createResponse.TradeIndex,
		Status:     types.StatusProcessed,
	}, *processResponse)
}
