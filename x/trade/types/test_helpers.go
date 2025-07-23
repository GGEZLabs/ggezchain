package types

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetSampleTradeDataJson(tradeType TradeType) string {
	tradeValue := 100.50
	tradeNetValue := 495.00
	numberOfShares := 1000.0
	sharePrice := 49.50
	shareNetPrice := 500.00

	var quantity *sdk.Coin
	if tradeType == TradeTypeBuy || tradeType == TradeTypeSell {
		quantity = &sdk.Coin{
			Denom:  DefaultDenom,
			Amount: math.NewInt(100000),
		}
	}

	switch tradeType {
	case TradeTypeSplit, TradeTypeReverseSplit:
		tradeValue = 0
		tradeNetValue = 0
		sharePrice = 0
		shareNetPrice = 0

	case TradeTypeDividends, TradeTypeDividendsDeduction:
		numberOfShares = 0
		sharePrice = 0
		shareNetPrice = 0
	}

	tradeData := TradeData{
		TradeInfo: &TradeInfo{
			AssetHolderId:       1,
			AssetId:             1,
			TradeType:           tradeType,
			TradeValue:          tradeValue,
			BaseCurrency:        "USD",
			SettlementCurrency:  "USD",
			ExchangeRate:        1,
			Exchange:            "US",
			FundName:            "TechFund",
			Issuer:              "CompanyA",
			NumberOfShares:      numberOfShares,
			CoinMintingPriceUsd: 0.001,
			Quantity:            quantity,
			Segment:             "Technology",
			SharePrice:          sharePrice,
			Ticker:              "TECH",
			TradeFee:            5.00,
			ShareNetPrice:       shareNetPrice,
			TradeNetValue:       tradeNetValue,
		},
		Brokerage: &Brokerage{
			Name:    "XYZBrokerage",
			Type:    "Online",
			Country: "USA",
		},
	}

	tradeDataBytes, err := json.Marshal(tradeData)
	if err != nil {
		panic(err)
	}
	return string(tradeDataBytes)
}

func GetSampleTradeData(tradeType TradeType) TradeData {
	tradeValue := 100.50
	tradeNetValue := 495.00
	numberOfShares := 1000.0
	sharePrice := 49.50
	shareNetPrice := 500.00

	var quantity *sdk.Coin
	if tradeType == TradeTypeBuy || tradeType == TradeTypeSell {
		quantity = &sdk.Coin{
			Denom:  DefaultDenom,
			Amount: math.NewInt(100000),
		}
	}

	switch tradeType {
	case TradeTypeSplit, TradeTypeReverseSplit:
		tradeValue = 0
		tradeNetValue = 0
		sharePrice = 0
		shareNetPrice = 0

	case TradeTypeDividends, TradeTypeDividendsDeduction:
		numberOfShares = 0
		sharePrice = 0
		shareNetPrice = 0
	}

	tradeData := TradeData{
		TradeInfo: &TradeInfo{
			AssetHolderId:       2,
			AssetId:             789,
			TradeType:           tradeType,
			TradeValue:          tradeValue,
			BaseCurrency:        "USD",
			SettlementCurrency:  "USD",
			ExchangeRate:        1,
			Exchange:            "NYSE",
			FundName:            "TechFund",
			Issuer:              "CompanyA",
			NumberOfShares:      numberOfShares,
			CoinMintingPriceUsd: 0.001,
			Quantity:            quantity,
			Segment:             "Technology",
			SharePrice:          sharePrice,
			Ticker:              "TECH",
			TradeFee:            5.00,
			ShareNetPrice:       shareNetPrice,
			TradeNetValue:       tradeNetValue,
		},
		Brokerage: &Brokerage{
			Name:    "XYZBrokerage",
			Type:    "Online",
			Country: "USA",
		},
	}
	return tradeData
}

func GetSampleExchangeRateJson() string {
	exchangeRate := []ExchangeRateJson{
		{
			FromCurrency:    "USD",
			ToCurrency:      "EUR",
			OriginalAmount:  1,
			ConvertedAmount: 0.85,
			CurrencyRate:    0.85,
			Timestamp:       "2025-07-08T00:00:00Z",
		},
	}

	exchangeRateBytes, err := json.Marshal(exchangeRate)
	if err != nil {
		panic(err)
	}
	return string(exchangeRateBytes)
}

func GetSampleCoinMintingPriceJson() string {
	coinMintingPrice := []CoinMintingPriceJson{
		{
			CurrencyCode: "USD",
			MintingPrice: 0.001,
		},
	}

	coinMintingPriceBytes, err := json.Marshal(coinMintingPrice)
	if err != nil {
		panic(err)
	}
	return string(coinMintingPriceBytes)
}

func GetBaseStoredTrade() StoredTrade {
	return StoredTrade{
		TradeType: TradeTypeBuy,
		Amount: &sdk.Coin{
			Denom:  DefaultDenom,
			Amount: math.NewInt(100000),
		},
		CoinMintingPriceUsd:  "0.001",
		ReceiverAddress:      testutil.Alice,
		Maker:                testutil.Alice,
		TradeData:            GetSampleTradeDataJson(TradeTypeBuy),
		BankingSystemData:    "{}",
		ExchangeRateJson:     GetSampleExchangeRateJson(),
		CoinMintingPriceJson: GetSampleCoinMintingPriceJson(),
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
		GetSampleTradeDataJson(TradeTypeBuy),
		"{}",
		GetSampleCoinMintingPriceJson(),
		GetSampleExchangeRateJson(),
	)
}

// GetMsgCreateTradeWithTypeAndAmount get sample create trade message specified with trade type and amount
func GetMsgCreateTradeWithTypeAndAmount(tradeType TradeType, amount int64) *MsgCreateTrade {
	tdStr := GetSampleTradeDataJson(TradeTypeBuy)
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
		Creator:              testutil.Alice,
		ReceiverAddress:      testutil.Alice,
		TradeData:            string(tdBytes),
		BankingSystemData:    "{}",
		CoinMintingPriceJson: GetSampleCoinMintingPriceJson(),
		ExchangeRateJson:     GetSampleExchangeRateJson(),
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
