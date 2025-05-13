package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestValidateTradeData(t *testing.T) {
	tests := []struct {
		name        string
		tradeData   string
		expectedErr bool
		errMsg      string
	}{
		{
			name:      "Valid trade data object",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
		},
		{
			name:      "Invalid trade data object",
			tradeData: `"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      types.ErrInvalidTradeDataObject.Error(),
		},
		{
			name:        "Invalid asset_holder_id",
			tradeData:   `{"trade_info":{"asset_holder_id":0,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "asset_holder_id must be greater than 0",
		},
		{
			name:        "Invalid asset_id",
			tradeData:   `{"trade_info":{"asset_holder_id":10,"asset_id":0,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "asset_id must be greater than 0",
		},
		{
			name:        "Invalid trade_value",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":0,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "trade_value must be greater than 0",
		},
		{
			name:        "Invalid currency",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "currency must not be empty or whitespace",
		},
		{
			name:        "Invalid exchange",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "exchange must not be empty or whitespace",
		},
		{
			name:        "Invalid fund_name",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "fund_name must not be empty or whitespace",
		},
		{
			name:        "Invalid issuer",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":" ","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "issuer must not be empty or whitespace",
		},
		{
			name:        "Invalid no_shares",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":0,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "no_shares must be greater than 0",
		},
		{
			name:        "Invalid price",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "price must be greater than 0",
		},
		{
			name:        "Invalid quantity",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":0,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "quantity must be greater than 0",
		},
		{
			name:        "Invalid segment",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":" ","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "segment must not be empty or whitespace",
		},
		{
			name:        "Invalid share_price",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":-5,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "share_price must be greater than 0",
		},
		{
			name:        "Invalid ticker",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":" ","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "ticker must not be empty or whitespace",
		},
		{
			name:        "Invalid trade_fee",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":-5,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "trade_fee must be a non-negative number",
		},
		{
			name:        "Invalid trade_net_price",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "trade_net_price must be greater than 0",
		},
		{
			name:        "Invalid trade_type",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":0,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "trade_type must be BUY or SELL",
		},
		{
			name:        "Invalid trade_net_value",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "trade_net_value must be greater than 0",
		},
		{
			name:        "Invalid brokerage country",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":""}}`,
			expectedErr: true,
			errMsg:      "brokerage country must not be empty or whitespace",
		},
		{
			name:        "Invalid brokerage type",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":" ","country":"US"}}`,
			expectedErr: true,
			errMsg:      "brokerage type must not be empty or whitespace",
		},
		{
			name:        "Invalid brokerage name",
			tradeData:   `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"","type":"Brokerage Firm","country":"US"}}`,
			expectedErr: true,
			errMsg:      "brokerage name must not be empty or whitespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateTradeData(tt.tradeData)
			if tt.expectedErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
