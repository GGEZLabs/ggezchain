package types_test

import (
	"testing"
	"time"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
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
