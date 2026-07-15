package types

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetSampleTradeDataJson returns a sample trade data JSON string for the given trade type.
// It mirrors GetSampleTradeData below but marshaled to JSON, useful for message payloads.
func GetSampleTradeDataJson(tradeType TradeType) string {
	tradeDataBytes, err := json.Marshal(GetSampleTradeData(tradeType))
	if err != nil {
		panic(err)
	}
	return string(tradeDataBytes)
}

// GetSampleTradeData returns a sample TradeData for the given trade type.
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

	return TradeData{
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
}

// GetSampleExchangeRateJson returns a sample exchange rate JSON payload.
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

// GetSampleCoinMintingPriceJson returns a sample coin minting price JSON payload.
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

// GetBaseStoredTrade returns the base sample StoredTrade shared by
// GetSampleStoredTrade/GetSampleStoredTradeConfirmed/GetSampleStoredTradeRejected below.
//
// NOTE: unlike the old repo, StoredTrade.Amount is a non-nullable sdk.Coin.
func GetBaseStoredTrade() StoredTrade {
	return StoredTrade{
		TradeType: TradeTypeBuy,
		Amount: sdk.Coin{
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

// GetSampleMsgCreateTrade returns a sample MsgCreateTrade.
//
// NOTE: unlike the old repo there is no NewMsgCreateTrade constructor in the
// new scaffold (message_create_trade.go was not ported; the equivalent
// validation now lives inline in the keeper's CreateTrade handler), so this
// builds the message struct literal directly.
func GetSampleMsgCreateTrade() *MsgCreateTrade {
	return &MsgCreateTrade{
		Creator:              testutil.Alice,
		ReceiverAddress:      testutil.Alice,
		TradeData:            GetSampleTradeDataJson(TradeTypeBuy),
		BankingSystemData:    "{}",
		CoinMintingPriceJson: GetSampleCoinMintingPriceJson(),
		ExchangeRateJson:     GetSampleExchangeRateJson(),
	}
}

// GetMsgCreateTradeWithTypeAndAmount returns a sample MsgCreateTrade for a buy/sell trade
// type with the given trade type and quantity amount.
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

// GetSampleStoredTrade returns the sample StoredTrade expected right after
// GetSampleMsgCreateTrade is created (status pending).
func GetSampleStoredTrade(tradeIndex uint64) StoredTrade {
	sampleStoredTrade := GetBaseStoredTrade()
	sampleStoredTrade.TradeIndex = tradeIndex
	sampleStoredTrade.Status = StatusPending
	sampleStoredTrade.Result = TradeCreatedSuccessfully

	return sampleStoredTrade
}

// GetSampleStoredTradeConfirmed returns the sample StoredTrade expected after a confirm process.
func GetSampleStoredTradeConfirmed(tradeIndex uint64) StoredTrade {
	sampleStoredTrade := GetBaseStoredTrade()
	sampleStoredTrade.TradeIndex = tradeIndex
	sampleStoredTrade.Status = StatusProcessed
	sampleStoredTrade.Checker = testutil.Bob
	sampleStoredTrade.Result = TradeProcessedSuccessfully

	return sampleStoredTrade
}

// GetSampleStoredTradeRejected returns the sample StoredTrade expected after a reject process.
func GetSampleStoredTradeRejected(tradeIndex uint64) StoredTrade {
	sampleStoredTrade := GetBaseStoredTrade()
	sampleStoredTrade.TradeIndex = tradeIndex
	sampleStoredTrade.Status = StatusRejected
	sampleStoredTrade.Checker = testutil.Bob
	sampleStoredTrade.Result = TradeProcessedSuccessfully

	return sampleStoredTrade
}
