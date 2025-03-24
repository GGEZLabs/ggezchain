package keeper_test

import (
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

func (suite *IntegrationTestSuite) TestCreateTrade() {
	suite.SetupTestForCreateTrade()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	suite.Nil(err)
	suite.EqualValues(types.MsgCreateTradeResponse{
		TradeIndex: 1,
		Status:     types.Pending,
	}, *createResponse)
}

func (suite *IntegrationTestSuite) TestIfTradeSaved() {
	suite.SetupTestForCreateTrade()
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)
	trade, found := keeper.GetStoredTrade(suite.ctx, 1)
	suite.True(found)
	suite.EqualValues(types.StoredTrade{
		TradeIndex:           1,
		Status:               types.Pending,
		CreateDate:           trade.CreateDate,
		TradeType:            types.Buy,
		Coin:                 types.DefaultCoinDenom,
		Price:                "0.001",
		Quantity:             "100000",
		ReceiverAddress:      testutil.Mutaz,
		Maker:                testutil.Mutaz,
		Checker:              "",
		ProcessDate:          trade.CreateDate,
		UpdateDate:           trade.CreateDate,
		TradeData:            "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		Result:               types.ErrTradeCreatedSuccessfully.Error(),
		BankingSystemData:    "{}",
		CoinMintingPriceJSON: "",
		ExchangeRateJSON:     ""}, trade)
}

func (suite *IntegrationTestSuite) TestIfTempTradeSaved() {
	suite.SetupTestForCreateTrade()
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)
	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, 1)
	suite.True(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     1,
		CreateDate:     tempTrade.CreateDate,
		TempTradeIndex: 1,
	}, tempTrade)
}

func (suite *IntegrationTestSuite) TestGetAllStoredTrade() {
	suite.SetupTestForCreateTrade()
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)
	allTrades := keeper.GetAllStoredTrade(suite.ctx)
	suite.EqualValues(types.StoredTrade{
		TradeIndex:           1,
		Status:               types.Pending,
		CreateDate:           allTrades[0].CreateDate,
		TradeType:            types.Buy,
		Coin:                 types.DefaultCoinDenom,
		Price:                "0.001",
		Quantity:             "100000",
		ReceiverAddress:      testutil.Mutaz,
		Maker:                testutil.Mutaz,
		Checker:              "",
		ProcessDate:          allTrades[0].CreateDate,
		UpdateDate:           allTrades[0].CreateDate,
		TradeData:            "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		Result:               types.ErrTradeCreatedSuccessfully.Error(),
		BankingSystemData:    "{}",
		CoinMintingPriceJSON: "",
		ExchangeRateJSON:     ""}, allTrades[0])
}

func (suite *IntegrationTestSuite) TestGetAllStoredTempTrade() {
	suite.SetupTestForCreateTrade()
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 2,
	}, tradeIndex)
	allTempTrades := keeper.GetAllStoredTempTrade(suite.ctx)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     1,
		CreateDate:     allTempTrades[0].CreateDate,
		TempTradeIndex: 1,
	}, allTempTrades[0])
}

func (suite *IntegrationTestSuite) TestCreateTradeWithInvalidMakerPermission() {
	suite.SetupTestForCreateTrade()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mohd,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})

	suite.Nil(createResponse)
	suite.ErrorIs(err, types.ErrInvalidMakerPermission)
}

func (suite *IntegrationTestSuite) TestCreateTradeWithInvalidTradeData() {
	suite.SetupTestForCreateTrade()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":0,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})

	suite.Nil(createResponse)
	suite.ErrorIs(err, types.ErrInvalidTradeDataObject)
}

func (suite *IntegrationTestSuite) TestCreate2Trades() {
	suite.SetupTestForCreateTrade()
	keeper := suite.app.TradeKeeper

	suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})
	trade, found := keeper.GetStoredTrade(suite.ctx, 1)
	suite.True(found)
	suite.EqualValues(types.StoredTrade{
		TradeIndex:           1,
		Status:               types.Pending,
		CreateDate:           trade.CreateDate,
		TradeType:            types.Buy,
		Coin:                 types.DefaultCoinDenom,
		Price:                "0.001",
		Quantity:             "100000",
		ReceiverAddress:      testutil.Mutaz,
		Maker:                testutil.Mutaz,
		Checker:              "",
		ProcessDate:          trade.CreateDate,
		UpdateDate:           trade.CreateDate,
		TradeData:            "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		Result:               types.ErrTradeCreatedSuccessfully.Error(),
		BankingSystemData:    "{}",
		CoinMintingPriceJSON: "",
		ExchangeRateJSON:     "",
	}, trade)

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, 1)
	suite.True(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     1,
		CreateDate:     tempTrade.CreateDate,
		TempTradeIndex: 1,
	}, tempTrade)

	suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:           testutil.Mutaz,
		TradeType:         types.Buy,
		Coin:              types.DefaultCoinDenom,
		Price:             "0.001",
		Quantity:          "100000",
		ReceiverAddress:   testutil.Mutaz,
		TradeData:         "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		BankingSystemData: "{}",
	})

	trade, found = keeper.GetStoredTrade(suite.ctx, 2)
	suite.True(found)
	suite.EqualValues(types.StoredTrade{
		TradeIndex:           2,
		Status:               types.Pending,
		CreateDate:           trade.CreateDate,
		TradeType:            types.Buy,
		Coin:                 types.DefaultCoinDenom,
		Price:                "0.001",
		Quantity:             "100000",
		ReceiverAddress:      testutil.Mutaz,
		Maker:                testutil.Mutaz,
		Checker:              "",
		ProcessDate:          trade.CreateDate,
		UpdateDate:           trade.CreateDate,
		TradeData:            "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":2,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":1000,\"price\":50.25,\"quantity\":10,\"segment\":\"Technology\",\"sharePrice\":49.50,\"ticker\":\"TECH\",\"tradeFee\":5.00,\"tradeNetPrice\":500.00,\"tradeNetValue\":495.00},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
		Result:               types.ErrTradeCreatedSuccessfully.Error(),
		BankingSystemData:    "{}",
		CoinMintingPriceJSON: "",
		ExchangeRateJSON:     ""}, trade)

	tempTrade, found = keeper.GetStoredTempTrade(suite.ctx, 2)
	suite.True(found)
	suite.EqualValues(types.StoredTempTrade{
		TradeIndex:     2,
		CreateDate:     tempTrade.CreateDate,
		TempTradeIndex: 2,
	}, tempTrade)

	// check get all trades and temp trades and next trade index
	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.True(found)
	suite.EqualValues(types.TradeIndex{
		NextId: 3,
	}, tradeIndex)
	AllTrades := keeper.GetAllStoredTrade(suite.ctx)
	suite.EqualValues(len(AllTrades), 2)

	AllTempTrades := keeper.GetAllStoredTempTrade(suite.ctx)
	suite.EqualValues(len(AllTempTrades), 2)
}
