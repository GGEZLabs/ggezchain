package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestIsTypeValid(t *testing.T) {
	tests := []struct {
		name      string
		tradeType types.TradeType
		valid     bool
	}{
		{
			name:      "trade type buy",
			tradeType: types.TradeTypeBuy,
			valid:     true,
		},
		{
			name:      "trade type sell",
			tradeType: types.TradeTypeSell,
			valid:     true,
		},
		{
			name:      "trade type split",
			tradeType: types.TradeTypeSplit,
			valid:     true,
		},
		{
			name:      "trade type reinvestment",
			tradeType: types.TradeTypeReinvestment,
			valid:     true,
		},
		{
			name:      "trade type dividends",
			tradeType: types.TradeTypeDividends,
			valid:     false,
		},
		{
			name:      "trade type unspecified",
			tradeType: types.TradeTypeUnspecified,
			valid:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.tradeType.IsTypeValid()
			require.Equal(t, tt.valid, valid)
		})
	}
}

func TestIsStatusValid(t *testing.T) {
	tests := []struct {
		name      string
		tradeType types.TradeStatus
		valid     bool
	}{
		{
			name:      "trade status processed",
			tradeType: types.StatusProcessed,
			valid:     true,
		},
		{
			name:      "trade status rejected",
			tradeType: types.StatusRejected,
			valid:     true,
		},
		{
			name:      "trade status failed",
			tradeType: types.StatusFailed,
			valid:     true,
		},
		{
			name:      "trade status canceled",
			tradeType: types.StatusCanceled,
			valid:     true,
		},
		{
			name:      "trade status pending",
			tradeType: types.StatusPending,
			valid:     true,
		},
		{
			name:      "trade status nil",
			tradeType: types.StatusNil,
			valid:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.tradeType.IsStatusValid()
			require.Equal(t, tt.valid, valid)
		})
	}
}
