package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestIsAddressWhiteListed(t *testing.T) {
	keeper, _ := keepertest.TradeKeeper(t)

	tests := []struct {
		name                 string
		address              string
		msgType              string
		isAddressWhitelisted bool
		err                  error
		ACLfilePath          string
	}{
		{
			name:                 "address white listed with type create trade",
			address:              testutil.Mutaz,
			msgType:              types.CreateTrade,
			err:                  nil,
			isAddressWhitelisted: true,
			ACLfilePath:          "/.ggezchain/config/chain_acl.json",
		},
		{
			name:                 "address is not white listed with type create trade",
			address:              "xxxx1ktj48fem0vw7kmg2gnny4dplmmvzvvmdatjjhk",
			msgType:              types.CreateTrade,
			err:                  nil,
			isAddressWhitelisted: false,
			ACLfilePath:          "/.ggezchain/config/chain_acl.json",
		}, {
			name:                 "address white listed with type process trade",
			address:              testutil.Mutaz,
			msgType:              types.ProcessTrade,
			err:                  nil,
			isAddressWhitelisted: true,
			ACLfilePath:          "/.ggezchain/config/chain_acl.json",
		},
		{
			name:                 "address is not white listed with type process trade",
			address:              "xxxx1ktj48fem0vw7kmg2gnny4dplmmvzvvmdatjjhk",
			msgType:              types.ProcessTrade,
			err:                  nil,
			isAddressWhitelisted: false,
			ACLfilePath:          "/.ggezchain/config/chain_acl.json",
		},
		{
			name:                 "address white listed with invalid type",
			address:              "ggez16rmzj02qkqp48v7rhlm99fvqrheenqepfeymsu",
			msgType:              "Invalid Type",
			err:                  nil,
			isAddressWhitelisted: false,
			ACLfilePath:          "/.ggezchain/config/chain_acl.json",
		},
		{
			name:                 "Invalid ACL file path",
			address:              "ggez16rmzj02qkqp48v7rhlm99fvqrheenqepfeymsu",
			msgType:              "Invalid Type",
			err:                  types.ErrInvalidPath,
			isAddressWhitelisted: false,
			ACLfilePath:          "invalid path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isAddressWhitelisted, err := keeper.IsAddressWhitelisted(tt.address, tt.msgType, tt.ACLfilePath)
			if err != nil {
				require.Equal(t, err, tt.err)
				return
			}
			require.Equal(t, isAddressWhitelisted, tt.isAddressWhitelisted)
		})
	}
}

func TestCancelExpiredPendingTrades(t *testing.T) {
	keeper, ctx := keepertest.TradeKeeper(t)
	now := ctx.BlockTime().UTC()
	formattedCurrentDate := now.Format("2006-01-02 15:04")

	storedTradeOne := types.StoredTrade{
		TradeIndex:      1,
		TradeType:       types.Buy,
		Coin:            "ugz",
		Price:           "10.59",
		Quantity:        "5026505",
		ReceiverAddress: sample.AccAddress(),
		TradeData:       "TradeData",
		Status:          "Pending",
		Maker:           sample.AccAddress(),
		Checker:         sample.AccAddress(),
		CreateDate:      "2023-05-07 12:00",
		UpdateDate:      "",
		ProcessDate:     "2023-05-07 12:00",
		Result:          "trade created successfully",
	}
	storedTradeTwo := types.StoredTrade{
		TradeIndex:      1,
		TradeType:       types.Sell,
		Coin:            "ugz",
		Price:           "10.32",
		Quantity:        "1000000005",
		ReceiverAddress: sample.AccAddress(),
		TradeData:       "TradeData",
		Status:          "Pending",
		Maker:           sample.AccAddress(),
		Checker:         sample.AccAddress(),
		CreateDate:      formattedCurrentDate,
		UpdateDate:      "",
		ProcessDate:     formattedCurrentDate,
		Result:          "trade created successfully",
	}
	storedTradeThree := types.StoredTrade{
		TradeIndex:      1,
		TradeType:       types.Buy,
		Coin:            "ugz",
		Price:           "15.19",
		Quantity:        "5000026000",
		ReceiverAddress: sample.AccAddress(),
		TradeData:       "TradeData",
		Status:          "Pending",
		Maker:           sample.AccAddress(),
		Checker:         sample.AccAddress(),
		CreateDate:      "2023-05-07 12:00",
		UpdateDate:      "",
		ProcessDate:     "",
		Result:          "trade created successfully",
	}

	// trade will be canceled
	storedTempTradeOne := types.StoredTempTrade{
		TradeIndex:     1,
		TempTradeIndex: 1,
		CreateDate:     "2023-05-07 12:00",
	}
	//trade not expired
	storedTempTradeTwo := types.StoredTempTrade{
		TradeIndex:     1,
		TempTradeIndex: 1,
		CreateDate:     formattedCurrentDate,
	}
	//invalid date
	storedTempTradeThree := types.StoredTempTrade{
		TradeIndex:     1,
		TempTradeIndex: 1,
		CreateDate:     "2023-05-06",
	}

	type TestFunc func()
	tests := []struct {
		name                     string
		expectedStatus           string
		expectedTempStoredLength int
		function                 TestFunc
		err                      error
	}{
		{
			name:                     "Check status for expired trade",
			expectedStatus:           types.Canceled,
			expectedTempStoredLength: 0,
			function: func() {
				keeper.SetStoredTrade(ctx, storedTradeOne)
				keeper.SetStoredTempTrade(ctx, storedTempTradeOne)
			},
			err: nil,
		},
		{
			name:                     "check status for not expired trade ",
			expectedStatus:           types.Pending,
			expectedTempStoredLength: 1,
			function: func() {
				keeper.SetStoredTrade(ctx, storedTradeTwo)
				keeper.SetStoredTempTrade(ctx, storedTempTradeTwo)
			},
			err: nil,
		},
		{
			name:                     "invalid date error",
			expectedStatus:           types.Pending,
			expectedTempStoredLength: 1,
			function: func() {
				keeper.SetStoredTrade(ctx, storedTradeThree)
				keeper.SetStoredTempTrade(ctx, storedTempTradeThree)
			},
			err: types.ErrInvalidDateFormat,
		},
	}

	for _, tt := range tests {
		//remove old temp trade
		allStoredTempTrade := keeper.GetAllStoredTempTrade(ctx)
		if len(allStoredTempTrade) > 1 {
			keeper.RemoveStoredTempTrade(ctx, 1)
			keeper.RemoveStoredTrade(ctx, 1)
		}
		tt.function()
		t.Run(tt.name, func(t *testing.T) {
			err := keeper.CancelExpiredPendingTrades(ctx)
			trade, _ := keeper.GetStoredTrade(ctx, 1)
			tempTradeLength := len(keeper.GetAllStoredTempTrade(ctx))
			if err != nil {
				require.ErrorIs(t, err, tt.err)
			}
			require.Equal(t, tt.expectedStatus, trade.Status)
			require.Equal(t, tt.expectedTempStoredLength, tempTradeLength)
		})
	}
}

func TestValidateTradeData(t *testing.T) {
	keeper, _ := keepertest.TradeKeeper(t)

	tests := []struct {
		name         string
		tradeDataObj string
		err          error
	}{
		{
			name:         "Valid trade data object",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          nil,
		},
		{
			name:         "Invalid trade data object",
			tradeDataObj: "\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrInvalidTradeDataJSON,
		},
		{
			name:         "Invalid assetHolderID",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":0,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataAssetHolderID,
		},
		{
			name:         "Invalid assetID",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":-6,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataAssetID,
		},
		{
			name:         "Invalid tradeRequestID",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":-1,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataRequestID,
		},
		{
			name:         "Invalid tradeValue",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":0,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataValue,
		},
		{
			name:         "Invalid currency",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataCurrency,
		},
		{
			name:         "Invalid exchange",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataExchange,
		},
		{
			name:         "Invalid fundName",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataFundName,
		},
		{
			name:         "Invalid issuer",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataIssuer,
		},
		{
			name:         "Invalid NoShares",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataNoShares,
		},
		{
			name:         "Invalid price",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataPrice,
		},
		{
			name:         "Invalid quantity",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataQuantity,
		},
		{
			name:         "Invalid segment",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataSegment,
		},
		{
			name:         "Invalid sharePrice",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataSharePrice,
		},
		{
			name:         "Invalid ticker",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataTicker,
		},
		{
			name:         "Invalid fee",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataFee,
		}, {
			name:         "Invalid netPrice",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrTradeDataNetPrice,
		},
		{
			name:         "Invalid tradeType",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrInvalidTradeType,
		},
		{
			name:         "Invalid brokerageCountry",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"Online\",\"country\":\"\"}}",
			err:          types.ErrBrokerageCountry,
		},
		{
			name:         "Invalid brokerageType",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"XYZBrokerage\",\"type\":\"\",\"country\":\"USA\"}}",
			err:          types.ErrBrokerageType,
		},
		{
			name:         "Invalid brokerageName",
			tradeDataObj: "{\"TradeData\":{\"tradeRequestID\":123,\"assetHolderID\":456,\"assetID\":789,\"tradeType\":\"Buy\",\"tradeValue\":100.50,\"currency\":\"USD\",\"exchange\":\"NYSE\",\"fundName\":\"TechFund\",\"issuer\":\"CompanyA\",\"noShares\":\"1000\",\"price\":\"50.25\",\"quantity\":\"10\",\"segment\":\"Technology\",\"sharePrice\":\"49.50\",\"ticker\":\"TECH\",\"tradeFee\":\"5.00\",\"tradeNetPrice\":\"500.00\",\"tradeNetValue\":\"495.00\"},\"Brokerage\":{\"name\":\"\",\"type\":\"Online\",\"country\":\"USA\"}}",
			err:          types.ErrBrokerageName,
		},
		{
			name:         "Invalid TradeDataJSON",
			tradeDataObj: "invalid obj",
			err:          types.ErrInvalidTradeDataJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := keeper.ValidateTradeData(tt.tradeDataObj)
			if err != nil {
				require.ErrorIs(t, err, tt.err)

				return
			}
		})
	}
}

func TestIsJSON(t *testing.T) {
	keeper, _ := keepertest.TradeKeeper(t)

	tests := []struct {
		name   string
		JSON   string
		isJSON bool
	}{
		{
			name:   "string is valid JSON",
			JSON:   "{\"browsers\":{\"firefox\":{\"name\":\"Firefox\",\"pref_url\":\"about:config\",\"releases\":{\"1\":{\"release_date\":\"2004-11-09\",\"status\":\"retired\",\"engine\":\"Gecko\",\"engine_version\":\"1.7\"}}}}}",
			isJSON: true,
		},
		{
			name:   "string is invalid JSON",
			JSON:   "Invalid JSON",
			isJSON: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isJSON := keeper.IsJSON(tt.JSON)

			require.Equal(t, isJSON, tt.isJSON)
		})
	}
}

func (suite *IntegrationTestSuite) TestMintOrBurnCoins() {
	suite.SetupTestForProcessTrade()
	var msgProcessTrade *types.MsgProcessTrade
	bankKeeper := suite.app.BankKeeper

	tests := []struct {
		name           string
		tradeData      types.StoredTrade
		expectedStatus string
		expectedSupply sdk.Coin
		err            error
	}{
		{
			name: "mint max amount coins",
			tradeData: types.StoredTrade{
				Coin:            types.DefaultCoinDenom,
				Quantity:        "9223372036854775807",
				ReceiverAddress: testutil.Mutaz,
				TradeType:       types.Buy,
			},
			expectedStatus: types.Completed,
			err:            types.ErrTradeProcessedSuccessfully,
			expectedSupply: sdk.Coin{
				Denom:  "ugz",
				Amount: sdkmath.NewInt(9223372036854775807),
			},
		},
		{
			name: "burn coins",
			tradeData: types.StoredTrade{
				Coin:            types.DefaultCoinDenom,
				Quantity:        "10000000",
				ReceiverAddress: testutil.Mutaz,
				TradeType:       types.Sell,
			},
			expectedStatus: types.Completed,
			err:            types.ErrTradeProcessedSuccessfully,
			expectedSupply: sdk.Coin{
				Denom:  "ugz",
				Amount: sdkmath.NewInt(9223372036844775807),
			},
		},
		{
			name: "mint coins",
			tradeData: types.StoredTrade{
				Coin:            types.DefaultCoinDenom,
				Quantity:        "10000000",
				ReceiverAddress: testutil.Mutaz,
				TradeType:       types.Buy,
			},
			expectedStatus: types.Completed,
			err:            types.ErrTradeProcessedSuccessfully,
			expectedSupply: sdk.Coin{
				Denom:  "ugz",
				Amount: sdkmath.NewInt(9223372036854775807),
			},
		},

		{
			name: "burn max amount coins",
			tradeData: types.StoredTrade{
				Coin:            types.DefaultCoinDenom,
				Quantity:        "9223372036854775807",
				ReceiverAddress: testutil.Mutaz,
				TradeType:       types.Sell,
			},
			expectedStatus: types.Completed,
			err:            types.ErrTradeProcessedSuccessfully,
			expectedSupply: sdk.Coin{
				Denom:  "ugz",
				Amount: sdkmath.NewInt(0),
			},
		},
		{
			name: "invalid receiver address",
			tradeData: types.StoredTrade{
				Coin:            types.DefaultCoinDenom,
				Quantity:        "10000000",
				ReceiverAddress: "invalidAddress",
				TradeType:       types.Sell,
			},
			expectedStatus: types.Failed,
			err:            types.ErrInvalidReceiverAddress,
			expectedSupply: sdk.Coin{
				Denom:  "ugz",
				Amount: sdkmath.NewInt(0),
			},
		},
		{
			name: "invalid quantity",
			tradeData: types.StoredTrade{
				Coin:            types.DefaultCoinDenom,
				Quantity:        "abc",
				ReceiverAddress: testutil.Rami,
				TradeType:       types.Sell,
			},
			expectedStatus: "",
			err:            types.ErrInvalidTradeQuantity,
			expectedSupply: sdk.Coin{
				Denom:  "ugz",
				Amount: sdkmath.NewInt(0),
			},
		},
		{
			name: "invalid quantity / exceeded max quantity",
			tradeData: types.StoredTrade{
				Coin:            types.DefaultCoinDenom,
				Quantity:        "92233720368547758077",
				ReceiverAddress: testutil.Ahmad,
				TradeType:       types.Sell,
			},
			expectedStatus: types.Completed,
			err:            types.ErrInvalidTradeQuantity,
			expectedSupply: sdk.Coin{
				Denom:  "ugz",
				Amount: sdkmath.NewInt(0),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(fmt.Sprintf("Case %s", tt.name), func() {
			supply := bankKeeper.GetSupply(suite.ctx, types.DefaultCoinDenom)

			coin, err := msgProcessTrade.GetPrepareCoin(tt.tradeData)
			if err != nil {
				suite.Equal(supply.Denom, tt.expectedSupply.Denom)
				suite.ErrorContains(err, tt.err.Error())
				return
			}
			status, err := suite.app.TradeKeeper.MintOrBurnCoins(suite.ctx, tt.tradeData, coin)
			if err != nil {
				suite.Equal(status, tt.expectedStatus)
				suite.ErrorIs(err, tt.err)
			}
			supply = bankKeeper.GetSupply(suite.ctx, types.DefaultCoinDenom)
			// check supply
			suite.Equal(supply, tt.expectedSupply)

		})
	}
}

func (suite *IntegrationTestSuite) TestIsAddressAllowed() {
	suite.SetupTestForProcessTrade()
	tests := []struct {
		name      string
		address   string
		msgType   string
		isAllowed bool
		err       error
	}{

		{
			name:      "address allowed to create trade",
			address:   testutil.Mutaz,
			msgType:   types.CreateTrade,
			isAllowed: true,
			err:       nil,
		},
		{
			name:      "address allowed to process trade",
			address:   testutil.Mohd,
			msgType:   types.ProcessTrade,
			isAllowed: true,
			err:       nil,
		},
		{
			name:      "address Not exist in ACL file",
			address:   testutil.Rami,
			msgType:   types.CreateTrade,
			isAllowed: false,
			err:       nil,
		},
		{
			name:      "address exist in ACL file but not delegate with any validator",
			address:   testutil.Rami,
			msgType:   types.CreateTrade,
			isAllowed: false,
			err:       nil,
		},
	}

	for _, tt := range tests {
		suite.Run(fmt.Sprintf("Case %s", tt.name), func() {
			isAllowed, err := suite.app.TradeKeeper.IsAddressAllowed(suite.ctx, tt.address, tt.msgType)
			if err != nil {
				suite.Equal(isAllowed, tt.isAllowed)
				suite.ErrorIs(err, tt.err)
				return
			}
			suite.Nil(err)
			suite.Equal(isAllowed, tt.isAllowed)
		})
	}
}
