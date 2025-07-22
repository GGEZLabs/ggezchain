package types

import (
	"encoding/json"
	"strings"
)

func ValidateTradeData(tradeData string) (TradeData, error) {
	var td TradeData
	if err := json.Unmarshal([]byte(tradeData), &td); err != nil {
		return td, ErrInvalidTradeData.Wrap(err.Error())
	}

	if td.TradeInfo == nil || td.Brokerage == nil {
		return td, ErrInvalidTradeData
	}

	if err := ValidateCommonTradeData(td); err != nil {
		return td, err
	}

	switch td.TradeInfo.TradeType {
	case TradeTypeBuy, TradeTypeSell:
		return td, ValidateBuyOrSell(td.TradeInfo)
	case TradeTypeReinvestment:
		return td, ValidateReinvestment(td.TradeInfo)
	case TradeTypeDividends, TradeTypeDividendsDeduction:
		return td, ValidateDividends(td.TradeInfo)
	case TradeTypeSplit, TradeTypeReverseSplit:
		return td, ValidateSplit(td.TradeInfo)
	default:
		return td, ErrInvalidTradeInfo.Wrapf("invalid trade_type, %s", td.TradeInfo.TradeType.String())
	}
}

// ValidateCommonTradeData validates fields common to all trade types
func ValidateCommonTradeData(td TradeData) error {
	if td.TradeInfo.AssetHolderId <= 0 {
		return ErrInvalidTradeInfo.Wrapf("asset_holder_id must be greater than 0, got: %d", td.TradeInfo.AssetHolderId)
	}
	if td.TradeInfo.AssetId <= 0 {
		return ErrInvalidTradeInfo.Wrapf("asset_id must be greater than 0, got: %d", td.TradeInfo.AssetId)
	}
	if !td.TradeInfo.TradeType.IsTypeValid() {
		return ErrInvalidTradeInfo.Wrapf("invalid trade_type, %s", td.TradeInfo.TradeType.String())
	}
	if strings.TrimSpace(td.TradeInfo.BaseCurrency) == "" {
		return ErrInvalidTradeInfo.Wrap("base_currency must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.SettlementCurrency) == "" {
		return ErrInvalidTradeInfo.Wrap("settlement_currency must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.Exchange) == "" {
		return ErrInvalidTradeInfo.Wrap("exchange must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.FundName) == "" {
		return ErrInvalidTradeInfo.Wrap("fund_name must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.Issuer) == "" {
		return ErrInvalidTradeInfo.Wrap("issuer must not be empty or whitespace")
	}
	if td.TradeInfo.CoinMintingPriceUsd <= 0 {
		return ErrInvalidTradeInfo.Wrapf("coin_minting_price_usd must be greater than 0, got: %f", td.TradeInfo.CoinMintingPriceUsd)
	}
	if strings.TrimSpace(td.TradeInfo.Segment) == "" {
		return ErrInvalidTradeInfo.Wrap("segment must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.Ticker) == "" {
		return ErrInvalidTradeInfo.Wrap("ticker must not be empty or whitespace")
	}
	if td.TradeInfo.TradeFee < 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_fee must be a non-negative number, got: %f", td.TradeInfo.TradeFee)
	}
	if td.TradeInfo.ExchangeRate <= 0 {
		return ErrInvalidTradeInfo.Wrapf("exchange_rate must be greater than 0, got: %f", td.TradeInfo.ExchangeRate)
	}
	if strings.TrimSpace(td.Brokerage.Country) == "" {
		return ErrInvalidTradeBrokerage.Wrap("brokerage country must not be empty or whitespace")
	}
	if strings.TrimSpace(td.Brokerage.Type) == "" {
		return ErrInvalidTradeBrokerage.Wrap("brokerage type must not be empty or whitespace")
	}
	if strings.TrimSpace(td.Brokerage.Name) == "" {
		return ErrInvalidTradeBrokerage.Wrap("brokerage name must not be empty or whitespace")
	}
	return nil
}

// ValidateBuyOrSell validates buy and sell trade types
func ValidateBuyOrSell(tradeInfo *TradeInfo) error {
	if tradeInfo.SharePrice <= 0 {
		return ErrInvalidTradeInfo.Wrapf("share_price must be greater than 0, got: %f", tradeInfo.SharePrice)
	}
	if tradeInfo.ShareNetPrice <= 0 {
		return ErrInvalidTradeInfo.Wrapf("share_net_price must be greater than 0, got: %f", tradeInfo.ShareNetPrice)
	}
	if tradeInfo.NumberOfShares <= 0 {
		return ErrInvalidTradeInfo.Wrapf("number_of_shares must be greater than 0, got: %f", tradeInfo.NumberOfShares)
	}
	if tradeInfo.TradeValue <= 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_value must be greater than 0, got: %f", tradeInfo.TradeValue)
	}
	if tradeInfo.TradeNetValue <= 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_net_value must be greater than 0, got: %f", tradeInfo.TradeNetValue)
	}
	if tradeInfo.Quantity == nil {
		return ErrInvalidTradeInfo.Wrap("invalid quantity")
	}
	if !tradeInfo.Quantity.IsValid() {
		return ErrInvalidTradeInfo.Wrapf("invalid quantity: %s", tradeInfo.Quantity.String())
	}
	if tradeInfo.Quantity.IsZero() {
		return ErrInvalidTradeInfo.Wrapf("zero quantity not allowed: %s", tradeInfo.Quantity.String())
	}
	if tradeInfo.Quantity.Denom != DefaultDenom {
		return ErrInvalidTradeInfo.Wrapf("invalid denom expected: %s, got: %s", DefaultDenom, tradeInfo.Quantity.Denom)
	}
	return nil
}

// ValidateReinvestment validates reinvestment trade type
func ValidateReinvestment(tradeInfo *TradeInfo) error {
	if tradeInfo.SharePrice <= 0 {
		return ErrInvalidTradeInfo.Wrapf("share_price must be greater than 0, got: %f", tradeInfo.SharePrice)
	}
	if tradeInfo.ShareNetPrice <= 0 {
		return ErrInvalidTradeInfo.Wrapf("share_net_price must be greater than 0, got: %f", tradeInfo.ShareNetPrice)
	}
	if tradeInfo.NumberOfShares <= 0 {
		return ErrInvalidTradeInfo.Wrapf("number_of_shares must be greater than 0, got: %f", tradeInfo.NumberOfShares)
	}
	if tradeInfo.TradeValue <= 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_value must be greater than 0, got: %f", tradeInfo.TradeValue)
	}
	if tradeInfo.TradeNetValue <= 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_net_value must be greater than 0, got: %f", tradeInfo.TradeNetValue)
	}
	if err := ValidateNoQuantity(tradeInfo); err != nil {
		return err
	}
	return nil
}

// ValidateDividends validates dividends and dividends deduction trade types
func ValidateDividends(tradeInfo *TradeInfo) error {
	if tradeInfo.SharePrice != 0 {
		return ErrInvalidTradeInfo.Wrapf("share_price must be 0, got: %f", tradeInfo.SharePrice)
	}
	if tradeInfo.ShareNetPrice != 0 {
		return ErrInvalidTradeInfo.Wrapf("share_net_price must be 0, got: %f", tradeInfo.ShareNetPrice)
	}
	if tradeInfo.NumberOfShares != 0 {
		return ErrInvalidTradeInfo.Wrapf("number_of_shares must be 0, got: %f", tradeInfo.NumberOfShares)
	}
	if tradeInfo.TradeValue <= 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_value must be greater than 0, got: %f", tradeInfo.TradeValue)
	}
	if tradeInfo.TradeNetValue <= 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_net_value must be greater than 0, got: %f", tradeInfo.TradeNetValue)
	}
	if err := ValidateNoQuantity(tradeInfo); err != nil {
		return err
	}
	return nil
}

// ValidateSplit validates split and reverse split trade types
func ValidateSplit(tradeInfo *TradeInfo) error {
	if tradeInfo.SharePrice != 0 {
		return ErrInvalidTradeInfo.Wrapf("share_price must be 0, got: %f", tradeInfo.SharePrice)
	}
	if tradeInfo.ShareNetPrice != 0 {
		return ErrInvalidTradeInfo.Wrapf("share_net_price must be 0, got: %f", tradeInfo.ShareNetPrice)
	}
	if tradeInfo.NumberOfShares <= 0 {
		return ErrInvalidTradeInfo.Wrapf("number_of_shares must be greater than 0, got: %f", tradeInfo.NumberOfShares)
	}
	if tradeInfo.TradeValue != 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_value must be 0, got: %f", tradeInfo.TradeValue)
	}
	if tradeInfo.TradeNetValue != 0 {
		return ErrInvalidTradeInfo.Wrapf("trade_net_value must be 0, got: %f", tradeInfo.TradeNetValue)
	}
	if err := ValidateNoQuantity(tradeInfo); err != nil {
		return err
	}
	return nil
}

// ValidateNoQuantity ensures no quantity is set for certain trade types
func ValidateNoQuantity(tradeInfo *TradeInfo) error {
	if tradeInfo.Quantity != nil &&
		(tradeInfo.Quantity.IsValid() || tradeInfo.Quantity.Denom != "" || !tradeInfo.Quantity.IsZero()) {
		return ErrInvalidTradeInfo.Wrapf("quantity must not be set for trade type %s, got: %s", tradeInfo.TradeType.String(), tradeInfo.Quantity.String())
	}
	return nil
}
