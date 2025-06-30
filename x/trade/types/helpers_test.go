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
