package types_test

import (
	"testing"
	"time"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestValidateDate(t *testing.T) {
	blockTime := time.Now()
	year := time.Now().Year() + 5
	futureDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	tests := []struct {
		name      string
		dateStr   string
		expErr    bool
		expErrMsg string
	}{
		{
			name:    "Valid date (before blockTime)",
			dateStr: "2025-06-11T23:59:59Z",
			expErr:  false,
		},
		{
			name:      "Invalid date (future date)",
			dateStr:   futureDate,
			expErr:    true,
			expErrMsg: "cannot be in the future",
		},
		{
			name:      "Invalid date format",
			dateStr:   "2023-05-11",
			expErr:    true,
			expErrMsg: "invalid date format",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateDate(blockTime, tt.dateStr)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1.2e-11, "0.000000000012"},
		{5.4e-11, "0.000000000054"},
		{0.000000000001, "0.000000000001"},
		{0.0001, "0.0001"},
		{123.456, "123.456"},
		{123.456000, "123.456"},
		{123.000000000000, "123"},
		{1e-13, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			price := types.FormatPrice(tt.input)
			require.Equal(t, tt.expected, price)
		})
	}
}

func TestValidateExchangeRateJson(t *testing.T) {
	tests := []struct {
		name             string
		exchangeRateJson string
		expErr           bool
		expErrMsg        string
	}{
		{
			name:             "all good",
			exchangeRateJson: types.GetSampleExchangeRateJson(),
			expErr:           false,
		},
		{
			name:             "Invalid json format",
			exchangeRateJson: `{"from_currency":"USD","to_currency":"EUR","original_amount":1,"converted_amount":0.85,"currency_rate":0.85,"timestamp":"2025-07-08T10:59:59Z"}`,
			expErr:           true,
			expErrMsg:        "invalid exchange rate json format",
		},
		{
			name:             "Invalid from_currency",
			exchangeRateJson: `[{"from_currency":"","to_currency":"EUR","original_amount":1,"converted_amount":0.85,"currency_rate":0.85,"timestamp":"2025-07-08T10:59:59Z"}]`,
			expErr:           true,
			expErrMsg:        "from_currency must not be empty or whitespace at index: 0",
		},
		{
			name:             "Invalid to_currency",
			exchangeRateJson: `[{"from_currency":"USD","to_currency":"","original_amount":1,"converted_amount":0.85,"currency_rate":0.85,"timestamp":"2025-07-08T10:59:59Z"}]`,
			expErr:           true,
			expErrMsg:        "to_currency must not be empty or whitespace at index: 0",
		},
		{
			name:             "Invalid original_amount",
			exchangeRateJson: `[{"from_currency":"USD","to_currency":"EUR","original_amount":-1,"converted_amount":0.85,"currency_rate":0.85,"timestamp":"2025-07-08T10:59:59Z"}]`,
			expErr:           true,
			expErrMsg:        "original_amount must be greater than 0",
		},
		{
			name:             "Invalid converted_amount",
			exchangeRateJson: `[{"from_currency":"USD","to_currency":"EUR","original_amount":1,"converted_amount":-1,"currency_rate":0.85,"timestamp":"2025-07-08T10:59:59Z"}]`,
			expErr:           true,
			expErrMsg:        "converted_amount must be greater than 0",
		},
		{
			name:             "Invalid currency_rate",
			exchangeRateJson: `[{"from_currency":"USD","to_currency":"EUR","original_amount":1,"converted_amount":0.85,"currency_rate":0,"timestamp":"2025-07-08T10:59:59Z"}]`,
			expErr:           true,
			expErrMsg:        "currency_rate must be greater than 0",
		},
		{
			name:             "Invalid timestamp",
			exchangeRateJson: `[{"from_currency":"USD","to_currency":"EUR","original_amount":1,"converted_amount":0.85,"currency_rate":0.85,"timestamp":"2025-07-0810:59:59Z"}]`,
			expErr:           true,
			expErrMsg:        "invalid timestamp format at index: 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateExchangeRateJson(tt.exchangeRateJson)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCoinMintingPriceJson(t *testing.T) {
	tests := []struct {
		name                 string
		coinMintingPriceJson string
		expErr               bool
		expErrMsg            string
	}{
		{
			name:                 "all good",
			coinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
			expErr:               false,
		},
		{
			name:                 "Invalid json format",
			coinMintingPriceJson: `{"currency_code":"USD","minting_price":1}`,
			expErr:               true,
			expErrMsg:            "invalid coin minting price json format",
		},
		{
			name:                 "Invalid currency_code",
			coinMintingPriceJson: `[{"currency_code":"","minting_price":1}]`,
			expErr:               true,
			expErrMsg:            "currency_code must not be empty or whitespace at index: 0",
		},
		{
			name:                 "Invalid minting_price",
			coinMintingPriceJson: `[{"currency_code":"USD","minting_price":0}]`,
			expErr:               true,
			expErrMsg:            "minting_price must be greater than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateCoinMintingPriceJson(tt.coinMintingPriceJson)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
