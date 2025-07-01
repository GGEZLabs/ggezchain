package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestValidateTradeData(t *testing.T) {
	tests := []struct {
		name      string
		tradeData string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "valid trade data object",
			tradeData: types.GetSampleTradeData(types.TradeTypeBuy),
		},
		{
			name:      "nil trade info",
			tradeData: `{"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: types.ErrInvalidTradeData.Error(),
		},
		{
			name:      "nil brokerage",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9}}`,
			expErr:    true,
			expErrMsg: types.ErrInvalidTradeData.Error(),
		},
		{
			name:      "invalid trade data object",
			tradeData: `"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: types.ErrInvalidTradeData.Error(),
		},
		{
			name:      "invalid asset_holder_id",
			tradeData: `{"trade_info":{"asset_holder_id":0,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "asset_holder_id must be greater than 0",
		},
		{
			name:      "invalid asset_id",
			tradeData: `{"trade_info":{"asset_holder_id":10,"asset_id":0,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "asset_id must be greater than 0",
		},
		{
			name:      "invalid trade_value",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":0,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_value must be greater than 0",
		},
		{
			name:      "invalid currency",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "currency must not be empty or whitespace",
		},
		{
			name:      "invalid exchange",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "exchange must not be empty or whitespace",
		},
		{
			name:      "invalid fund_name",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "fund_name must not be empty or whitespace",
		},
		{
			name:      "invalid issuer",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":" ","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "issuer must not be empty or whitespace",
		},
		{
			name:      "invalid no_shares",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":0,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "no_shares must be greater than 0",
		},
		{
			name:      "invalid price",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "price must be greater than 0",
		},
		{
			name:      "invalid quantity",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"0","denom":"ug"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid quantity",
		},
		{
			name:      "zero quantity",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"0","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "zero quantity not allowed",
		},
		{
			name:      "invalid denom",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggez1"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid denom",
		},
		{
			name:      "send quantity with trade type split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid quantity: quantity must be zero",
		},
		{
			name:      "send quantity with trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":4,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid quantity: quantity must be zero",
		},
		{
			name:      "invalid segment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":" ","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "segment must not be empty or whitespace",
		},
		{
			name:      "invalid share_price",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":-5,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_price must be greater than 0",
		},
		{
			name:      "invalid ticker",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":" ","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "ticker must not be empty or whitespace",
		},
		{
			name:      "invalid trade_fee",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":-5,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_fee must be a non-negative number",
		},
		{
			name:      "invalid trade_net_price",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_net_price must be greater than 0",
		},
		{
			name:      "invalid trade_type",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":0,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid trade_type",
		},
		{
			name:      "invalid trade_net_value",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_net_value must be greater than 0",
		},
		{
			name:      "invalid brokerage country",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":""}}`,
			expErr:    true,
			expErrMsg: "brokerage country must not be empty or whitespace",
		},
		{
			name:      "invalid brokerage type",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":" ","country":"US"}}`,
			expErr:    true,
			expErrMsg: "brokerage type must not be empty or whitespace",
		},
		{
			name:      "invalid brokerage name",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "brokerage name must not be empty or whitespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := types.ValidateTradeData(tt.tradeData)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
