package keeper_test

import (
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestProcessTradeConfirm() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	processResponse, err := suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mohd,
		ProcessType: types.Confirm,
		TradeIndex:  1,
	})
	suite.Nil(err)
	suite.EqualValues(types.MsgProcessTradeResponse{
		TradeIndex:  1,
		Status:      types.Completed,
		Checker:     "",
		Maker:       "",
		TradeData:   "",
		CreateDate:  "",
		UpdateDate:  "",
		ProcessDate: "",
	}, *processResponse)
}

func (suite *IntegrationTestSuite) TestProcessTradeReject() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	processResponse, err := suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mohd,
		ProcessType: types.Reject,
		TradeIndex:  1,
	})
	suite.Nil(err)
	suite.EqualValues(types.MsgProcessTradeResponse{
		TradeIndex:  1,
		Status:      types.Rejected,
		Checker:     "",
		Maker:       "",
		TradeData:   "",
		CreateDate:  "",
		UpdateDate:  "",
		ProcessDate: "",
	}, *processResponse)
}

func (suite *IntegrationTestSuite) TestProcessTradeWithInvalidCheckerPermission() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	_, err := suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Rami,
		ProcessType: types.Confirm,
		TradeIndex:  1,
	})
	suite.ErrorIs(err, types.ErrInvalidCheckerPermission)
}

func (suite *IntegrationTestSuite) TestStoredTradeAfterConfirmTrade() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mohd,
		ProcessType: types.Confirm,
		TradeIndex:  1,
	})

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)

	trade, found := keeper.GetStoredTrade(suite.ctx, 1)

	suite.True(found)
	suite.EqualValues(types.StoredTrade{
		TradeIndex:      1,
		Status:          types.Completed,
		CreateDate:      trade.CreateDate,
		TradeType:       types.Buy,
		Coin:            types.DefaultCoinDenom,
		Price:           "0.001",
		Quantity:        "100000",
		ReceiverAddress: testutil.Mutaz,
		Maker:           testutil.Mutaz,
		Checker:         testutil.Mohd,
		ProcessDate:     trade.CreateDate,
		TradeData:       "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		Result:          types.ErrTradeProcessedSuccessfully.Error(),
		UpdateDate:      trade.UpdateDate,
	}, trade)
}

func (suite *IntegrationTestSuite) TestTempTradeAfterConfirmTrade() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mohd,
		ProcessType: types.Confirm,
		TradeIndex:  1,
	})

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, 1)

	suite.False(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     0,
		CreateDate:     "",
		TempTradeIndex: 0,
	}, tempTrade)
}

func (suite *IntegrationTestSuite) TestStoredTradeAfterRejectTrade() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mohd,
		ProcessType: types.Reject,
		TradeIndex:  1,
	})

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)

	trade, found := keeper.GetStoredTrade(suite.ctx, 1)

	suite.True(found)
	suite.EqualValues(types.StoredTrade{
		TradeIndex:      1,
		Status:          types.Rejected,
		CreateDate:      trade.CreateDate,
		TradeType:       types.Buy,
		Coin:            types.DefaultCoinDenom,
		Price:           "0.001",
		Quantity:        "100000",
		ReceiverAddress: testutil.Mutaz,
		Maker:           testutil.Mutaz,
		Checker:         testutil.Mohd,
		ProcessDate:     trade.CreateDate,
		TradeData:       "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		Result:          types.ErrTradeProcessedSuccessfully.Error(),
		UpdateDate:      trade.UpdateDate,
	}, trade)
}

func (suite *IntegrationTestSuite) TestTempTradeAfterRejectTrade() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mohd,
		ProcessType: types.Reject,
		TradeIndex:  1,
	})

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, 1)

	suite.False(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     0,
		CreateDate:     "",
		TempTradeIndex: 0,
	}, tempTrade)
}

func (suite *IntegrationTestSuite) TestProcessTwoTrade() {
	suite.SetupTestForProcessTrade()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})

	suite.msgServer.CreateTrade(goCtx, &types.MsgCreateTrade{
		Creator:           testutil.Mohd,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mohd,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)

	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 3,
	}, tradeIndex)

	tempTrades := keeper.GetAllStoredTempTrade(suite.ctx)
	trades := keeper.GetAllStoredTrade(suite.ctx)
	suite.EqualValues(len(tempTrades), 2)
	suite.EqualValues(len(trades), 2)

	suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mohd,
		ProcessType: types.Confirm,
		TradeIndex:  1,
	})

	suite.msgServer.ProcessTrade(goCtx, &types.MsgProcessTrade{
		Creator:     testutil.Mutaz,
		ProcessType: types.Reject,
		TradeIndex:  2,
	})

	tempTrades = keeper.GetAllStoredTempTrade(suite.ctx)
	trades = keeper.GetAllStoredTrade(suite.ctx)
	suite.EqualValues(len(tempTrades), 0)
	suite.EqualValues(len(trades), 2)
}
