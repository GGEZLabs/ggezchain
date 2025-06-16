package types

import (
	"encoding/json"
	"strings"
)

func ValidateTradeData(tradeData string) (TradeData, error) {
	var td TradeData
	if err := json.Unmarshal([]byte(tradeData), &td); err != nil {
		return td, ErrInvalidTradeData.Wrapf(err.Error())
	}

	if td.TradeInfo == nil || td.Brokerage == nil {
		return td, ErrInvalidTradeData
	}

	if td.TradeInfo.AssetHolderId <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("asset_holder_id must be greater than 0, got: %d", td.TradeInfo.AssetHolderId)
	}
	if td.TradeInfo.AssetId <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("asset_id must be greater than 0, got: %d", td.TradeInfo.AssetId)
	}
	if td.TradeInfo.TradeValue <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("trade_value must be greater than 0, got: %f", td.TradeInfo.TradeValue)
	}
	if strings.TrimSpace(td.TradeInfo.Currency) == "" {
		return td, ErrInvalidTradeInfo.Wrapf("currency must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.Exchange) == "" {
		return td, ErrInvalidTradeInfo.Wrapf("exchange must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.FundName) == "" {
		return td, ErrInvalidTradeInfo.Wrapf("fund_name must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.Issuer) == "" {
		return td, ErrInvalidTradeInfo.Wrapf("issuer must not be empty or whitespace")
	}
	if td.TradeInfo.NoShares <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("no_shares must be greater than 0, got: %d", td.TradeInfo.NoShares)
	}
	if td.TradeInfo.Price <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("price must be greater than 0, got: %f", td.TradeInfo.Price)
	}
	if !td.TradeInfo.Quantity.IsValid() {
		return td, ErrInvalidTradeInfo.Wrapf("invalid quantity: %s", td.TradeInfo.Quantity)
	}
	if td.TradeInfo.Quantity.Amount.IsZero() {
		return td, ErrInvalidTradeInfo.Wrapf("zero quantity not allowed: %s", td.TradeInfo.Quantity)
	}
	if td.TradeInfo.Quantity.Denom != DefaultDenom {
		return td, ErrInvalidTradeInfo.Wrapf("invalid denom expected: %s, got: %s ", DefaultDenom, td.TradeInfo.Quantity.Denom)
	}
	if strings.TrimSpace(td.TradeInfo.Segment) == "" {
		return td, ErrInvalidTradeInfo.Wrapf("segment must not be empty or whitespace")
	}
	if td.TradeInfo.SharePrice <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("share_price must be greater than 0, got: %f", td.TradeInfo.SharePrice)
	}
	if strings.TrimSpace(td.TradeInfo.Ticker) == "" {
		return td, ErrInvalidTradeInfo.Wrapf("ticker must not be empty or whitespace")
	}
	if td.TradeInfo.TradeFee < 0 {
		return td, ErrInvalidTradeInfo.Wrapf("trade_fee must be a non-negative number, got: %f", td.TradeInfo.TradeFee)
	}
	if td.TradeInfo.TradeNetPrice <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("trade_net_price must be greater than 0, got: %f", td.TradeInfo.TradeNetPrice)
	}
	if td.TradeInfo.TradeNetValue <= 0 {
		return td, ErrInvalidTradeInfo.Wrapf("trade_net_value must be greater than 0, got: %f", td.TradeInfo.TradeNetValue)
	}
	if td.TradeInfo.TradeType != TradeTypeBuy &&
		td.TradeInfo.TradeType != TradeTypeSell {
		return td, ErrInvalidTradeInfo.Wrapf("trade_type must be BUY or SELL")
	}
	if strings.TrimSpace(td.Brokerage.Country) == "" {
		return td, ErrInvalidTradeBrokerage.Wrapf("brokerage country must not be empty or whitespace")
	}
	if strings.TrimSpace(td.Brokerage.Type) == "" {
		return td, ErrInvalidTradeBrokerage.Wrapf("brokerage type must not be empty or whitespace")
	}
	if strings.TrimSpace(td.Brokerage.Name) == "" {
		return td, ErrInvalidTradeBrokerage.Wrapf("brokerage name must not be empty or whitespace")
	}
	return td, nil
}
