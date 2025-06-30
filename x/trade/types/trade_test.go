package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestIsValid(t *testing.T) {
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
			valid:     true,
		},
		{
			name:      "trade type unspecified",
			tradeType: types.TradeTypeUnspecified,
			valid:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.tradeType.IsValid()
			require.Equal(t, tt.valid, valid)
		})
	}
}
