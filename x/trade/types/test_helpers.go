package types

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetSampleTradeData(tradeType TradeType) string {
	tradeData := TradeData{
		TradeInfo: &TradeInfo{
			AssetHolderId: 2,
			AssetId:       789,
			TradeType:     tradeType,
			TradeValue:    100.50,
			Currency:      "USD",
			Exchange:      "NYSE",
			FundName:      "TechFund",
			Issuer:        "CompanyA",
			NoShares:      1000,
			Price:         0.001,
			Quantity: &sdk.Coin{
				Denom:  DefaultDenom,
				Amount: math.NewInt(100000),
			},
			Segment:       "Technology",
			SharePrice:    49.50,
			Ticker:        "TECH",
			TradeFee:      5.00,
			TradeNetPrice: 500.00,
			TradeNetValue: 495.00,
		},
		Brokerage: &Brokerage{
			Name:    "XYZBrokerage",
			Type:    "Online",
			Country: "USA",
		},
	}

	tradeDataJSON, err := json.Marshal(tradeData)
	if err != nil {
		panic(err)
	}
	return string(tradeDataJSON)
}

func GetBaseStoredTrade() StoredTrade {
	return StoredTrade{
		TradeType: TradeTypeBuy,
		Amount: &sdk.Coin{
			Denom:  DefaultDenom,
			Amount: math.NewInt(100000),
		},
		Price:                "0.001",
		ReceiverAddress:      testutil.Alice,
		Maker:                testutil.Alice,
		TradeData:            GetSampleTradeData(TradeTypeBuy),
		BankingSystemData:    "{}",
		CoinMintingPriceJson: "",
		ExchangeRateJson:     "",
		TxDate:               "0001-01-01T00:00:00Z",
		CreateDate:           "0001-01-01T00:00:00Z",
		ProcessDate:          "0001-01-01T00:00:00Z",
		UpdateDate:           "0001-01-01T00:00:00Z",
	}
}

// GetSampleMsgCreateTrade get sample create trade message
func GetSampleMsgCreateTrade() *MsgCreateTrade {
	return NewMsgCreateTrade(
		testutil.Alice,
		testutil.Alice,
		GetSampleTradeData(TradeTypeBuy),
		"{}",
		"",
		"",
	)
}

// GetMsgCreateTradeWithTypeAndAmount get sample create trade message specified with trade type and amount
func GetMsgCreateTradeWithTypeAndAmount(tradeType TradeType, amount int64) *MsgCreateTrade {
	tdStr := GetSampleTradeData(TradeTypeBuy)
	var td TradeData
	if err := json.Unmarshal([]byte(tdStr), &td); err != nil {
		panic(err)
	}

	td.TradeInfo.TradeType = tradeType
	td.TradeInfo.Quantity.Amount = math.NewInt(amount)

	tdBytes, err := json.Marshal(td)
	if err != nil {
		panic(err)
	}

	return &MsgCreateTrade{
		Creator:           testutil.Alice,
		ReceiverAddress:   testutil.Alice,
		TradeData:         string(tdBytes),
		BankingSystemData: "{}",
	}
}

// GetSampleStoredTrade return sample stored trade according to GetSampleMsgCreateTrade function
// used after create trade
func GetSampleStoredTrade(tradeIndex uint64) StoredTrade {
	sampleStoredTrade := GetBaseStoredTrade()
	sampleStoredTrade.TradeIndex = tradeIndex
	sampleStoredTrade.Status = StatusPending
	sampleStoredTrade.Result = TradeCreatedSuccessfully

	return sampleStoredTrade
}

// GetSampleStoredTradeConfirmed used after confirm trade
func GetSampleStoredTradeConfirmed(tradeIndex uint64) StoredTrade {
	sampleStoredTrade := GetBaseStoredTrade()
	sampleStoredTrade.TradeIndex = tradeIndex
	sampleStoredTrade.Status = StatusProcessed
	sampleStoredTrade.Checker = testutil.Bob
	sampleStoredTrade.Result = TradeProcessedSuccessfully

	return sampleStoredTrade
}

// GetSampleStoredTradeRejected used after reject trade
func GetSampleStoredTradeRejected(tradeIndex uint64) StoredTrade {
	sampleStoredTrade := GetBaseStoredTrade()
	sampleStoredTrade.TradeIndex = tradeIndex
	sampleStoredTrade.Status = StatusRejected
	sampleStoredTrade.Checker = testutil.Bob
	sampleStoredTrade.Result = TradeProcessedSuccessfully

	return sampleStoredTrade
}
