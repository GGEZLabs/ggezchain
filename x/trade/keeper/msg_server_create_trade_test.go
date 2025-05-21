package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestCreateTrade() {
	indexes := suite.createTrade(1)
	suite.Require().EqualValues(1, len(indexes))
	suite.Require().EqualValues(1, indexes[0])
}

func (suite *KeeperTestSuite) TestIfTradeSaved() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found := keeper.GetStoredTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestIfTempTradeSaved() {
	indexes := suite.createTrade(1)
	keeper := suite.app.TradeKeeper

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTempTrade{
		TradeIndex:     indexes[0],
		CreateDate:     types.GetSampleStoredTrade(indexes[0]).CreateDate,
		TempTradeIndex: indexes[0],
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestGetAllStoredTrade() {
	indexes := suite.createTrade(3)
	keeper := suite.app.TradeKeeper

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	allTrades := keeper.GetAllStoredTrade(suite.ctx)
	suite.Require().EqualValues(len(allTrades), len(indexes))
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), allTrades[0])
}

func (suite *KeeperTestSuite) TestGetAllStoredTempTrade() {
	indexes := suite.createTrade(5)
	keeper := suite.app.TradeKeeper

	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	allTempTrades := keeper.GetAllStoredTempTrade(suite.ctx)
	suite.Require().EqualValues(types.StoredTempTrade{
		TradeIndex:     indexes[0],
		CreateDate:     types.GetSampleStoredTrade(indexes[0]).CreateDate,
		TempTradeIndex: indexes[0],
	}, allTempTrades[0])
}

func (suite *KeeperTestSuite) TestCreateTradeWithInvalidMakerPermission() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:   testutil.Bob,
		TradeType: types.TradeTypeBuy,
		Amount: &sdk.Coin{
			Denom:  types.DefaultCoinDenom,
			Amount: sdkmath.NewInt(100000),
		},
		Price:             "0.001",
		ReceiverAddress:   testutil.Alice,
		TradeData:         `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData: "{}",
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidMakerPermission)
}

func (suite *KeeperTestSuite) TestCreateTradeAuthorityAddressNotExist() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:   testutil.Eve,
		TradeType: types.TradeTypeBuy,
		Amount: &sdk.Coin{
			Denom:  types.DefaultCoinDenom,
			Amount: sdkmath.NewInt(100000),
		},
		Price:             "0.001",
		ReceiverAddress:   testutil.Alice,
		TradeData:         `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData: "{}",
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, acltypes.ErrAuthorityAddressNotExist)
	suite.Require().Contains(err.Error(), "unauthorized account")
}

func (suite *KeeperTestSuite) TestCreateTradeNoPermissionForModule() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:   testutil.Carol,
		TradeType: types.TradeTypeBuy,
		Amount: &sdk.Coin{
			Denom:  types.DefaultCoinDenom,
			Amount: sdkmath.NewInt(100000),
		},
		Price:             "0.001",
		ReceiverAddress:   testutil.Alice,
		TradeData:         `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData: "{}",
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrModuleNotFound)
	suite.Require().Contains(err.Error(), "no permission for module trade")
}

func (suite *KeeperTestSuite) TestCreateTradeWithInvalidTradeData() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:   testutil.Alice,
		TradeType: types.TradeTypeBuy,
		Amount: &sdk.Coin{
			Denom:  types.DefaultCoinDenom,
			Amount: sdkmath.NewInt(100000),
		},
		Price:             "0.001",
		ReceiverAddress:   testutil.Alice,
		TradeData:         `{"trade_info":{"asset_holder_id":0,"asset_id":1,"trade_type":0,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData: "{}",
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidTradeInfo)
}

func (suite *KeeperTestSuite) TestCreateTrades() {
	indexes := suite.createTrade(1000)
	keeper := suite.app.TradeKeeper

	for _, tradeIndex := range indexes {
		trade, found := keeper.GetStoredTrade(suite.ctx, tradeIndex)
		suite.Require().True(found)
		suite.Require().EqualValues(types.GetSampleStoredTrade(tradeIndex), trade)

		tempTrade, found := keeper.GetStoredTempTrade(suite.ctx, tradeIndex)
		suite.Require().True(found)
		suite.Require().EqualValues(types.StoredTempTrade{
			TradeIndex:     tradeIndex,
			CreateDate:     tempTrade.CreateDate,
			TempTradeIndex: tradeIndex,
		}, tempTrade)
	}

	// check get all trades and temp trades and next trade index
	tradeIndex, found := keeper.GetTradeIndex(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)
	AllTrades := keeper.GetAllStoredTrade(suite.ctx)
	suite.Require().EqualValues(len(AllTrades), uint64(len(indexes)))

	AllTempTrades := keeper.GetAllStoredTempTrade(suite.ctx)
	suite.Require().EqualValues(len(AllTempTrades), uint64(len(indexes)))
}
