package types

import (
	"encoding/json"
	"strings"
)

func ValidateTradeData(tradeData string) (err error) {
	var td TradeData
	if err := json.Unmarshal([]byte(tradeData), &td); err != nil {
		return ErrInvalidTradeDataObject
	}

	if td.TradeInfo.AssetHolderId <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("asset_holder_id must be greater than 0, got: %d", td.TradeInfo.AssetHolderId)
	}
	if td.TradeInfo.AssetId <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("asset_id must be greater than 0, got: %d", td.TradeInfo.AssetId)
	}
	if td.TradeInfo.TradeValue <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("trade_value must be greater than 0, got: %f", td.TradeInfo.TradeValue)
	}
	if strings.TrimSpace(td.TradeInfo.Currency) == "" {
		return ErrInvalidTradeDataObject.Wrapf("currency must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.Exchange) == "" {
		return ErrInvalidTradeDataObject.Wrapf("exchange must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.FundName) == "" {
		return ErrInvalidTradeDataObject.Wrapf("fund_name must not be empty or whitespace")
	}
	if strings.TrimSpace(td.TradeInfo.Issuer) == "" {
		return ErrInvalidTradeDataObject.Wrapf("issuer must not be empty or whitespace")
	}
	if td.TradeInfo.NoShares <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("no_shares must be greater than 0, got: %d", td.TradeInfo.NoShares)
	}
	if td.TradeInfo.Price <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("price must be greater than 0, got: %f", td.TradeInfo.Price)
	}
	if td.TradeInfo.Quantity <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("quantity must be greater than 0, got: %d", td.TradeInfo.Quantity)
	}
	if strings.TrimSpace(td.TradeInfo.Segment) == "" {
		return ErrInvalidTradeDataObject.Wrapf("segment must not be empty or whitespace")
	}
	if td.TradeInfo.SharePrice <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("share_price must be greater than 0, got: %f", td.TradeInfo.SharePrice)
	}
	if strings.TrimSpace(td.TradeInfo.Ticker) == "" {
		return ErrInvalidTradeDataObject.Wrapf("ticker must not be empty or whitespace")
	}
	if td.TradeInfo.TradeFee < 0 {
		return ErrInvalidTradeDataObject.Wrapf("trade_fee must be a non-negative number, got: %f", td.TradeInfo.TradeFee)
	}
	if td.TradeInfo.TradeNetPrice <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("trade_net_price must be greater than 0, got: %f", td.TradeInfo.TradeNetPrice)
	}
	if td.TradeInfo.TradeNetValue <= 0 {
		return ErrInvalidTradeDataObject.Wrapf("trade_net_value must be greater than 0, got: %f", td.TradeInfo.TradeNetValue)
	}
	if td.TradeInfo.TradeType != TradeTypeBuy &&
		td.TradeInfo.TradeType != TradeTypeSell {
		return ErrInvalidTradeDataObject.Wrapf("trade_type must be BUY or SELL")
	}
	if strings.TrimSpace(td.Brokerage.Country) == "" {
		return ErrInvalidTradeDataObject.Wrapf("brokerage country must not be empty or whitespace")
	}
	if strings.TrimSpace(td.Brokerage.Type) == "" {
		return ErrInvalidTradeDataObject.Wrapf("brokerage type must not be empty or whitespace")
	}
	if strings.TrimSpace(td.Brokerage.Name) == "" {
		return ErrInvalidTradeDataObject.Wrapf("brokerage name must not be empty or whitespace")
	}
	return nil

}
