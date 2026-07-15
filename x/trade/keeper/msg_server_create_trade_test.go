package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	tradetestutil "github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestCreateTrade() {
	indexes := suite.createNTrades(1)
	suite.Require().EqualValues(1, len(indexes))
	suite.Require().EqualValues(1, indexes[0])
}

func (suite *KeeperTestSuite) TestIfTradeSaved() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	trade, found := getStoredTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), trade)
}

func (suite *KeeperTestSuite) TestIfTempTradeSaved() {
	indexes := suite.createNTrades(1)
	keeper := suite.tradeKeeper

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	tempTrade, found := getStoredTempTrade(suite.ctx, keeper, indexes[0])
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTempTrade{
		TradeIndex: indexes[0],
		TxDate:     types.GetSampleStoredTrade(indexes[0]).TxDate,
	}, tempTrade)
}

func (suite *KeeperTestSuite) TestGetAllStoredTrade() {
	indexes := suite.createNTrades(3)
	keeper := suite.tradeKeeper

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	allTrades := getAllStoredTrade(suite.ctx, keeper)
	suite.Require().EqualValues(len(allTrades), len(indexes))
	suite.Require().EqualValues(types.GetSampleStoredTrade(indexes[0]), allTrades[0])
}

func (suite *KeeperTestSuite) TestGetAllStoredTempTrade() {
	indexes := suite.createNTrades(5)
	keeper := suite.tradeKeeper

	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)

	allTempTrades := getAllStoredTempTrade(suite.ctx, keeper)
	suite.Require().EqualValues(types.StoredTempTrade{
		TradeIndex: indexes[0],
		TxDate:     types.GetSampleStoredTrade(indexes[0]).TxDate,
	}, allTempTrades[0])
}

// TestCreateTradeValidation covers the former message_create_trade_test.go
// ValidateBasic() cases. That validation now lives inline in the keeper's
// CreateTrade handler (there's no more standalone MsgCreateTrade.ValidateBasic
// in the new scaffold), so these are ported as keeper-level cases instead.
func (suite *KeeperTestSuite) TestCreateTradeWithInvalidMakerPermission() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Bob,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidMakerPermission)
}

func (suite *KeeperTestSuite) TestCreateTradeAuthorityAddressNotExist() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Eve,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, acltypes.ErrAuthorityAddressDoesNotExist)
	suite.Require().Contains(err.Error(), "unauthorized account")
}

func (suite *KeeperTestSuite) TestCreateTradeNoPermissionForModule() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Carol, // Does not have permission for trade module
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrModuleNotFound)
	suite.Require().Contains(err.Error(), "no permission for module trade")
}

func (suite *KeeperTestSuite) TestCreateTradeWithInvalidTradeData() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            `{"trade_info":{"asset_holder_id":0,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidTradeInfo)
}

func (suite *KeeperTestSuite) TestCreateTradeWithCoinMintingPriceJson() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleExchangeRateJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidCoinMintingPriceJson)
	suite.Require().Contains(err.Error(), "currency_code must not be empty or whitespace at index: 0")
}

func (suite *KeeperTestSuite) TestCreateTradeWithExchangeRateJson() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleCoinMintingPriceJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidExchangeRateJson)
	suite.Require().Contains(err.Error(), "from_currency must not be empty or whitespace at index: 0")
}

func (suite *KeeperTestSuite) TestCreateTrades() {
	indexes := suite.createNTrades(1000)
	keeper := suite.tradeKeeper

	for _, tradeIndex := range indexes {
		trade, found := getStoredTrade(suite.ctx, keeper, tradeIndex)
		suite.Require().True(found)
		suite.Require().EqualValues(types.GetSampleStoredTrade(tradeIndex), trade)

		tempTrade, found := getStoredTempTrade(suite.ctx, keeper, tradeIndex)
		suite.Require().True(found)
		suite.Require().EqualValues(types.StoredTempTrade{
			TradeIndex: tradeIndex,
			TxDate:     tempTrade.TxDate,
		}, tempTrade)
	}

	// check get all trades and temp trades and next trade index
	tradeIndex, found := getTradeIndex(suite.ctx, keeper)
	suite.Require().True(found)
	suite.Require().EqualValues(types.TradeIndex{
		NextId: uint64(len(indexes) + 1),
	}, tradeIndex)
	allTrades := getAllStoredTrade(suite.ctx, keeper)
	suite.Require().EqualValues(len(allTrades), uint64(len(indexes)))

	allTempTrades := getAllStoredTempTrade(suite.ctx, keeper)
	suite.Require().EqualValues(len(allTempTrades), uint64(len(indexes)))
}

func (suite *KeeperTestSuite) TestCreateTradeWithInvalidCreateDate() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		CreateDate:           "2023-05-06",
	})

	suite.Require().Nil(createResponse)
	suite.Require().Contains(err.Error(), "invalid date format")
}

func (suite *KeeperTestSuite) TestCreateTradeWithCreateDateInFuture() {
	suite.setupTest()
	blockHeight := int64(1)
	blockTime := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockHeight(blockHeight).WithBlockTime(blockTime)

	futureDate := time.Date(blockTime.Year()+5, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	// Use EXPECT after updating context
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, tradetestutil.Alice).Return(acltypes.AclAuthority{
		Address: tradetestutil.Alice,
		Name:    "Alice",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	}, nil).AnyTimes()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		CreateDate:           futureDate,
	})

	suite.Require().Nil(createResponse)
	suite.Require().Contains(err.Error(), "date cannot be in the future")
}

func (suite *KeeperTestSuite) TestCreateTradeWithValidCreateDate() {
	suite.setupTest()
	blockHeight := int64(1)
	blockTime := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockHeight(blockHeight).WithBlockTime(blockTime)

	// Use EXPECT after updating context
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, tradetestutil.Alice).Return(acltypes.AclAuthority{
		Address: tradetestutil.Alice,
		Name:    "Alice",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	}, nil).AnyTimes()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		CreateDate:           "2024-05-11T08:44:00Z",
	})

	suite.Require().Nil(err)
	suite.Require().Equal(&types.MsgCreateTradeResponse{
		TradeIndex: uint64(1),
		Status:     types.StatusPending,
	}, createResponse)

	trade, found := getStoredTrade(suite.ctx, suite.tradeKeeper, 1)
	suite.Require().True(found)
	suite.Require().Equal(trade.CreateDate, "2024-05-11T08:44:00Z")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeSplit() {
	suite.setupTest()
	keeper := suite.tradeKeeper

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeSplit),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})
	suite.Require().Equal(&types.MsgCreateTradeResponse{
		TradeIndex: 1,
		Status:     types.StatusPending,
	}, createResponse)
	suite.Require().NoError(err)

	trade, found := getStoredTrade(suite.ctx, keeper, 1)
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTrade{
		TradeIndex:           1,
		TradeType:            types.TradeTypeSplit,
		Amount:               sdk.Coin{Amount: math.NewInt(0)},
		CoinMintingPriceUsd:  "0.001",
		Status:               types.StatusPending,
		Maker:                tradetestutil.Alice,
		TxDate:               "0001-01-01T00:00:00Z",
		CreateDate:           "0001-01-01T00:00:00Z",
		UpdateDate:           "0001-01-01T00:00:00Z",
		ProcessDate:          "0001-01-01T00:00:00Z",
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"TechFund","issuer":"CompanyA","number_of_shares":1000,"coin_minting_price_usd":0.001,"segment":"Technology","ticker":"TECH","trade_fee":5},"brokerage":{"name":"XYZBrokerage","type":"Online","country":"USA"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		Result:               types.TradeCreatedSuccessfully,
	}, trade)
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeSplitAndReceiverAddress() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeSplit),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().Contains(err.Error(), "receiver address must not be set for trade type TRADE_TYPE_SPLIT")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeSplitAndQuantity() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"trade_value":0,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidTradeInfo)
	suite.Require().Contains(err.Error(), "quantity must not be set")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeReverseSplit() {
	suite.setupTest()
	keeper := suite.tradeKeeper

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeReverseSplit),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Equal(&types.MsgCreateTradeResponse{
		TradeIndex: 1,
		Status:     types.StatusPending,
	}, createResponse)
	suite.Require().NoError(err)

	trade, found := getStoredTrade(suite.ctx, keeper, 1)
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTrade{
		TradeIndex:           1,
		TradeType:            types.TradeTypeReverseSplit,
		Amount:               sdk.Coin{Amount: math.NewInt(0)},
		CoinMintingPriceUsd:  "0.001",
		Status:               types.StatusPending,
		Maker:                tradetestutil.Alice,
		TxDate:               "0001-01-01T00:00:00Z",
		CreateDate:           "0001-01-01T00:00:00Z",
		UpdateDate:           "0001-01-01T00:00:00Z",
		ProcessDate:          "0001-01-01T00:00:00Z",
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":4,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"TechFund","issuer":"CompanyA","number_of_shares":1000,"coin_minting_price_usd":0.001,"segment":"Technology","ticker":"TECH","trade_fee":5},"brokerage":{"name":"XYZBrokerage","type":"Online","country":"USA"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		Result:               types.TradeCreatedSuccessfully,
	}, trade)
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeReverseSplitAndReceiverAddress() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeReverseSplit),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().Contains(err.Error(), "receiver address must not be set for trade type TRADE_TYPE_REVERSE_SPLIT")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeReverseSplitAndQuantity() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":4,"trade_value":0,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.001,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidTradeInfo)
	suite.Require().Contains(err.Error(), "quantity must not be set")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeReinvestment() {
	suite.setupTest()
	keeper := suite.tradeKeeper

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeReinvestment),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Equal(&types.MsgCreateTradeResponse{
		TradeIndex: 1,
		Status:     types.StatusPending,
	}, createResponse)
	suite.Require().NoError(err)

	trade, found := getStoredTrade(suite.ctx, keeper, 1)
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTrade{
		TradeIndex:           1,
		TradeType:            types.TradeTypeReinvestment,
		Amount:               sdk.Coin{Amount: math.NewInt(0)},
		CoinMintingPriceUsd:  "0.001",
		Status:               types.StatusPending,
		Maker:                tradetestutil.Alice,
		TxDate:               "0001-01-01T00:00:00Z",
		CreateDate:           "0001-01-01T00:00:00Z",
		UpdateDate:           "0001-01-01T00:00:00Z",
		ProcessDate:          "0001-01-01T00:00:00Z",
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":100.5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"TechFund","issuer":"CompanyA","number_of_shares":1000,"coin_minting_price_usd":0.001,"segment":"Technology","share_price":49.5,"ticker":"TECH","trade_fee":5,"share_net_price":500,"trade_net_value":495},"brokerage":{"name":"XYZBrokerage","type":"Online","country":"USA"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		Result:               types.TradeCreatedSuccessfully,
	}, trade)
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeReinvestmentAndReceiverAddress() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeReinvestment),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().Contains(err.Error(), "receiver address must not be set for trade type TRADE_TYPE_REINVESTMENT")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeReinvestmentAndQuantity() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.001,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidTradeInfo)
	suite.Require().Contains(err.Error(), "quantity must not be set")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeDividends() {
	suite.setupTest()
	keeper := suite.tradeKeeper

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeDividends),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Equal(&types.MsgCreateTradeResponse{
		TradeIndex: 1,
		Status:     types.StatusPending,
	}, createResponse)
	suite.Require().NoError(err)

	trade, found := getStoredTrade(suite.ctx, keeper, 1)
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTrade{
		TradeIndex:           1,
		TradeType:            types.TradeTypeDividends,
		Amount:               sdk.Coin{Amount: math.NewInt(0)},
		CoinMintingPriceUsd:  "0.001",
		Status:               types.StatusPending,
		Maker:                tradetestutil.Alice,
		TxDate:               "0001-01-01T00:00:00Z",
		CreateDate:           "0001-01-01T00:00:00Z",
		UpdateDate:           "0001-01-01T00:00:00Z",
		ProcessDate:          "0001-01-01T00:00:00Z",
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":6,"trade_value":100.5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"TechFund","issuer":"CompanyA","coin_minting_price_usd":0.001,"segment":"Technology","ticker":"TECH","trade_fee":5,"trade_net_value":495},"brokerage":{"name":"XYZBrokerage","type":"Online","country":"USA"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		Result:               types.TradeCreatedSuccessfully,
	}, trade)
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeDividendsAndReceiverAddress() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeDividends),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().Contains(err.Error(), "receiver address must not be set for trade type TRADE_TYPE_DIVIDENDS")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeDividendsAndQuantity() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":6,"trade_value":100.5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"TechFund","issuer":"CompanyA","coin_minting_price_usd":0.001,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Technology","ticker":"TECH","trade_fee":5,"trade_net_value":495},"brokerage":{"name":"XYZBrokerage","type":"Online","country":"USA"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidTradeInfo)
	suite.Require().Contains(err.Error(), "quantity must not be set")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeDividendsDeduction() {
	suite.setupTest()
	keeper := suite.tradeKeeper

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeDividendsDeduction),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Equal(&types.MsgCreateTradeResponse{
		TradeIndex: 1,
		Status:     types.StatusPending,
	}, createResponse)
	suite.Require().NoError(err)

	trade, found := getStoredTrade(suite.ctx, keeper, 1)
	suite.Require().True(found)
	suite.Require().EqualValues(types.StoredTrade{
		TradeIndex:           1,
		TradeType:            types.TradeTypeDividendsDeduction,
		Amount:               sdk.Coin{Amount: math.NewInt(0)},
		CoinMintingPriceUsd:  "0.001",
		Status:               types.StatusPending,
		Maker:                tradetestutil.Alice,
		TxDate:               "0001-01-01T00:00:00Z",
		CreateDate:           "0001-01-01T00:00:00Z",
		UpdateDate:           "0001-01-01T00:00:00Z",
		ProcessDate:          "0001-01-01T00:00:00Z",
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":7,"trade_value":100.5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"TechFund","issuer":"CompanyA","coin_minting_price_usd":0.001,"segment":"Technology","ticker":"TECH","trade_fee":5,"trade_net_value":495},"brokerage":{"name":"XYZBrokerage","type":"Online","country":"USA"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
		Result:               types.TradeCreatedSuccessfully,
	}, trade)
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeDividendsDeductionAndReceiverAddress() {
	suite.setupTest()

	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		ReceiverAddress:      tradetestutil.Alice,
		TradeData:            types.GetSampleTradeDataJson(types.TradeTypeDividendsDeduction),
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().Contains(err.Error(), "receiver address must not be set for trade type TRADE_TYPE_DIVIDEND_DEDUCTION")
}

func (suite *KeeperTestSuite) TestCreateTradeWithTypeDividendsDeductionAndQuantity() {
	suite.setupTest()
	createResponse, err := suite.msgServer.CreateTrade(suite.ctx, &types.MsgCreateTrade{
		Creator:              tradetestutil.Alice,
		TradeData:            `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":7,"trade_value":100.5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"TechFund","issuer":"CompanyA","coin_minting_price_usd":0.001,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Technology","ticker":"TECH","trade_fee":5,"trade_net_value":495},"brokerage":{"name":"XYZBrokerage","type":"Online","country":"USA"}}`,
		BankingSystemData:    "{}",
		ExchangeRateJson:     types.GetSampleExchangeRateJson(),
		CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
	})

	suite.Require().Nil(createResponse)
	suite.Require().ErrorIs(err, types.ErrInvalidTradeInfo)
	suite.Require().Contains(err.Error(), "quantity must not be set")
}
