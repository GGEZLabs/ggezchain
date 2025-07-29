package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateDate checks if the input date string is in the correct format and not after today's date.
func ValidateDate(blockTime time.Time, dateStr string) error {
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid date format: %s", err)
	}

	now := blockTime.Truncate(24 * time.Hour)
	parsedDate = parsedDate.Truncate(24 * time.Hour)

	if parsedDate.After(now) {
		return sdkerrors.ErrInvalidRequest.Wrapf("date cannot be in the future %s", parsedDate.Format(time.RFC3339))
	}

	return nil
}

// FormatPrice convert a float to a decimal string
func FormatPrice(price float64) string {
	str := fmt.Sprintf("%.12f", price)
	return strings.TrimRight(strings.TrimRight(str, "0"), ".")
}

// ValidateExchangeRateJson
func ValidateExchangeRateJson(exchangeRateJson string) error {
	var exchangeRates []ExchangeRateJson
	if err := json.Unmarshal([]byte(exchangeRateJson), &exchangeRates); err != nil {
		return ErrInvalidExchangeRateJson.Wrap(err.Error())
	}
	for i, rate := range exchangeRates {
		if strings.TrimSpace(rate.FromCurrency) == "" {
			return ErrInvalidExchangeRateJson.Wrapf("from_currency must not be empty or whitespace at index: %d", i)
		}
		if strings.TrimSpace(rate.ToCurrency) == "" {
			return ErrInvalidExchangeRateJson.Wrapf("to_currency must not be empty or whitespace at index: %d", i)
		}
		if rate.OriginalAmount < 0 {
			return ErrInvalidExchangeRateJson.Wrapf("original_amount must be a non-negative number, got: %f, at index: %d", rate.OriginalAmount, i)
		}
		if rate.ConvertedAmount < 0 {
			return ErrInvalidExchangeRateJson.Wrapf("converted_amount must be a non-negative number, got: %f, at index: %d", rate.ConvertedAmount, i)
		}
		if rate.CurrencyRate <= 0 {
			return ErrInvalidExchangeRateJson.Wrapf("currency_rate must be greater than 0, got: %f, at index: %d", rate.CurrencyRate, i)
		}
		_, err := time.Parse(time.RFC3339, rate.Timestamp)
		if err != nil {
			return ErrInvalidExchangeRateJson.Wrapf("invalid timestamp format at index: %d: %s", i, err)
		}
	}
	return nil
}

// ValidateCoinMintingPriceJson
func ValidateCoinMintingPriceJson(coinMintingPriceJson string) error {
	var coinMintingPrices []CoinMintingPriceJson
	if err := json.Unmarshal([]byte(coinMintingPriceJson), &coinMintingPrices); err != nil {
		return ErrInvalidCoinMintingPriceJson.Wrap(err.Error())
	}
	for i, mintingPrice := range coinMintingPrices {
		if strings.TrimSpace(mintingPrice.CurrencyCode) == "" {
			return ErrInvalidCoinMintingPriceJson.Wrapf("currency_code must not be empty or whitespace at index: %d", i)
		}
		if mintingPrice.MintingPrice <= 0 {
			return ErrInvalidCoinMintingPriceJson.Wrapf("minting_price must be greater than 0, got: %f, at index: %d", mintingPrice.MintingPrice, i)
		}
	}
	return nil
}
