package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
			tradeData: types.GetSampleTradeDataJson(types.TradeTypeBuy),
		},
		{
			name:      "nil trade info",
			tradeData: `{"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: types.ErrInvalidTradeData.Error(),
		},
		{
			name:      "nil brokerage",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9}}`,
			expErr:    true,
			expErrMsg: types.ErrInvalidTradeData.Error(),
		},
		{
			name:      "invalid trade data object",
			tradeData: `"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: types.ErrInvalidTradeData.Error(),
		},
		// Test common trade data
		{
			name:      "invalid asset_holder_id",
			tradeData: `{"trade_info":{"asset_holder_id":0,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "asset_holder_id must be greater than 0",
		},
		{
			name:      "invalid asset_id",
			tradeData: `{"trade_info":{"asset_holder_id":10,"asset_id":0,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "asset_id must be greater than 0",
		},
		{
			name:      "invalid trade_type",
			tradeData: `{"trade_info":{"asset_holder_id":10,"asset_id":1,"trade_type":0,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid trade_type",
		},
		{
			name:      "invalid base currency",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "base_currency must not be empty or whitespace",
		},
		{
			name:      "invalid settlement currency",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "settlement_currency must not be empty or whitespace",
		},
		{
			name:      "invalid exchange",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "exchange must not be empty or whitespace",
		},
		{
			name:      "invalid fund_name",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "fund_name must not be empty or whitespace",
		},
		{
			name:      "invalid issuer",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":" ","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "issuer must not be empty or whitespace",
		},
		{
			name:      "invalid price",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "coin_minting_price_usd must be greater than 0",
		},
		{
			name:      "invalid segment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":" ","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "segment must not be empty or whitespace",
		},
		{
			name:      "invalid ticker",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":" ","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "ticker must not be empty or whitespace",
		},
		{
			name:      "invalid trade_fee",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":-5,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_fee must be a non-negative number",
		},
		{
			name:      "invalid exchange rate",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":0,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "exchange_rate must be greater than 0, got: 0",
		},
		{
			name:      "invalid brokerage country",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":""}}`,
			expErr:    true,
			expErrMsg: "brokerage country must not be empty or whitespace",
		},
		{
			name:      "invalid brokerage type",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":" ","country":"US"}}`,
			expErr:    true,
			expErrMsg: "brokerage type must not be empty or whitespace",
		},
		{
			name:      "invalid brokerage name",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "brokerage name must not be empty or whitespace",
		},
		// Test buy and sell types
		{
			name:      "invalid share_price - trade type buy",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_price must be greater than 0",
		},
		{
			name:      "invalid share_net_price - trade type buy",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_net_price must be greater than 0",
		},
		{
			name:      "invalid number_of_shares - trade type buy",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "number_of_shares must be greater than 0",
		},
		{
			name:      "invalid trade value - trade type buy",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":-5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_value must be greater than 0",
		},
		{
			name:      "invalid trade net value - trade type sell",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":2,"trade_value":100,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_net_value must be greater than 0",
		},
		{
			name:      "invalid quantity - trade type sell",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":2,"trade_value":100,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid quantity",
		},
		{
			name:      "invalid quantity - trade type sell",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":2,"trade_value":100,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid quantity",
		},
		{
			name:      "zero quantity - trade type sell",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":2,"trade_value":100,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"0","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "zero quantity not allowed",
		},
		{
			name:      "invalid denom - trade type sell",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":2,"trade_value":100,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggez1"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "invalid denom expected: uggz, got: uggez1",
		},
		// Test reinvestment type
		{
			name:      "invalid share_price - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_price must be greater than 0",
		},
		{
			name:      "invalid share_net_price - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_net_price must be greater than 0",
		},
		{
			name:      "invalid number_of_shares - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "number_of_shares must be greater than 0",
		},
		{
			name:      "invalid trade value - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":-5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_value must be greater than 0",
		},
		{
			name:      "invalid trade net value - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_net_value must be greater than 0",
		},
		{
			name:      "invalid quantity - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "quantity must not be set for trade type TRADE_TYPE_REINVESTMENT",
		},
		{
			name:      "set valid quantity - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"0","denom":""},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    false,
		},
		{
			name:      "nil quantity - trade type reinvestment",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":5,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    false,
		},
		// Test dividends and dividends deduction types
		{
			name:      "invalid share_price - trade type dividends",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":6,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_price must be 0",
		},
		{
			name:      "invalid share_net_price - trade type dividends",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":6,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_net_price must be 0",
		},
		{
			name:      "invalid number_of_shares - trade type dividends",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":6,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "number_of_shares must be 0",
		},
		{
			name:      "invalid trade value - trade type dividends",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":6,"trade_value":-5,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_value must be greater than 0",
		},
		{
			name:      "invalid trade net value - trade type dividends deduction",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":7,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_net_value must be greater than 0",
		},
		{
			name:      "invalid quantity - trade type dividends deduction",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":7,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "quantity must not be set for trade type TRADE_TYPE_DIVIDEND_DEDUCTION",
		},
		{
			name:      "set valid quantity - trade type dividends deduction",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":7,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"0","denom":""},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    false,
		},
		{
			name:      "nil quantity - trade type dividends deduction",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":7,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    false,
		},
		// Test split and reverse split types
		{
			name:      "invalid share_price - trade type split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":100,"ticker":"CRBN","trade_fee":0,"share_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_price must be 0",
		},
		{
			name:      "invalid share_net_price - trade type split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":100,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "share_net_price must be 0",
		},
		{
			name:      "invalid number_of_shares - trade type split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":0,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "number_of_shares must be greater than 0",
		},
		{
			name:      "invalid trade value - trade type split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"trade_value":15,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_value must be 0",
		},
		{
			name:      "invalid trade net value - trade type reverse split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":4,"trade_value":0,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "trade_net_value must be 0",
		},
		{
			name:      "invalid quantity - trade type split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":3,"trade_value":0,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "quantity must not be set for trade type TRADE_TYPE_SPLIT",
		},
		{
			name:      "invalid quantity - trade type reverse split",
			tradeData: `{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":4,"trade_value":0,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","number_of_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":0,"ticker":"CRBN","trade_fee":0,"share_net_price":0,"trade_net_value":0},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			expErr:    true,
			expErrMsg: "quantity must not be set for trade type TRADE_TYPE_REVERSE_SPLIT",
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

func TestValidateCommonTradeData(t *testing.T) {
	tests := []struct {
		name      string
		tradeData types.TradeData
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "valid trade data object",
			tradeData: types.GetSampleTradeData(types.TradeTypeBuy),
		},
		{
			name: "invalid asset_holder_id",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId: 0,
				},
			},
			expErr:    true,
			expErrMsg: "asset_holder_id must be greater than 0",
		},
		{
			name: "invalid asset_id",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId: 1,
					AssetId:       0,
				},
			},
			expErr:    true,
			expErrMsg: "asset_id must be greater than 0",
		},
		{
			name: "invalid trade_type",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId: 1,
					AssetId:       1,
					TradeType:     types.TradeTypeUnspecified,
				},
			},
			expErr:    true,
			expErrMsg: "invalid trade_type",
		},
		{
			name: "invalid base currency",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId: 1,
					AssetId:       1,
					TradeType:     types.TradeTypeBuy,
					BaseCurrency:  " ",
				},
			},
			expErr:    true,
			expErrMsg: "base_currency must not be empty or whitespace",
		},
		{
			name: "invalid settlement currency",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:      1,
					AssetId:            1,
					TradeType:          types.TradeTypeBuy,
					BaseCurrency:       "USD",
					SettlementCurrency: "",
				},
			},
			expErr:    true,
			expErrMsg: "settlement_currency must not be empty or whitespace",
		},
		{
			name: "invalid exchange",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:      1,
					AssetId:            1,
					TradeType:          types.TradeTypeBuy,
					BaseCurrency:       "USD",
					SettlementCurrency: "USD",
					Exchange:           "",
				},
			},
			expErr:    true,
			expErrMsg: "exchange must not be empty or whitespace",
		},
		{
			name: "invalid fund_name",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:      1,
					AssetId:            1,
					TradeType:          types.TradeTypeBuy,
					BaseCurrency:       "USD",
					SettlementCurrency: "USD",
					Exchange:           "US",
					FundName:           "",
				},
			},
			expErr:    true,
			expErrMsg: "fund_name must not be empty or whitespace",
		},
		{
			name: "invalid issuer",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:      1,
					AssetId:            1,
					TradeType:          types.TradeTypeBuy,
					BaseCurrency:       "USD",
					SettlementCurrency: "USD",
					Exchange:           "US",
					FundName:           "TechFund",
					Issuer:             "",
				},
			},
			expErr:    true,
			expErrMsg: "issuer must not be empty or whitespace",
		},
		{
			name: "invalid price",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 0,
				},
			},
			expErr:    true,
			expErrMsg: "coin_minting_price_usd must be greater than 0",
		},
		{
			name: "invalid segment",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 1,
					Segment:             "",
				},
			},
			expErr:    true,
			expErrMsg: "segment must not be empty or whitespace",
		},
		{
			name: "invalid ticker",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 1,
					Segment:             "Global Low Carbon",
					Ticker:              "",
				},
			},
			expErr:    true,
			expErrMsg: "ticker must not be empty or whitespace",
		},
		{
			name: "invalid trade_fee",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 1,
					Segment:             "Global Low Carbon",
					Ticker:              "CEN",
					TradeFee:            -1,
				},
			},
			expErr:    true,
			expErrMsg: "trade_fee must be a non-negative number",
		},
		{
			name: "invalid exchange rate",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 1,
					Segment:             "Global Low Carbon",
					Ticker:              "CEN",
					TradeFee:            0,
					ExchangeRate:        0,
				},
			},
			expErr:    true,
			expErrMsg: "exchange_rate must be greater than 0",
		},
		{
			name: "invalid brokerage country",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 1,
					Segment:             "Global Low Carbon",
					Ticker:              "CEN",
					TradeFee:            0,
					ExchangeRate:        1,
				},
				Brokerage: &types.Brokerage{
					Country: "",
				},
			},
			expErr:    true,
			expErrMsg: "brokerage country must not be empty or whitespace",
		},
		{
			name: "invalid brokerage type",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 1,
					Segment:             "Global Low Carbon",
					Ticker:              "CEN",
					TradeFee:            0,
					ExchangeRate:        1,
				},
				Brokerage: &types.Brokerage{
					Country: "USA",
					Type:    "",
				},
			},
			expErr:    true,
			expErrMsg: "brokerage type must not be empty or whitespace",
		},
		{
			name: "invalid brokerage name",
			tradeData: types.TradeData{
				TradeInfo: &types.TradeInfo{
					AssetHolderId:       1,
					AssetId:             1,
					TradeType:           types.TradeTypeBuy,
					BaseCurrency:        "USD",
					SettlementCurrency:  "USD",
					Exchange:            "US",
					FundName:            "TechFund",
					Issuer:              "Blackrock",
					CoinMintingPriceUsd: 1,
					Segment:             "Global Low Carbon",
					Ticker:              "CEN",
					TradeFee:            0,
					ExchangeRate:        1,
				},
				Brokerage: &types.Brokerage{
					Country: "USA",
					Type:    "Brokerage Firm",
					Name:    "",
				},
			},
			expErr:    true,
			expErrMsg: "brokerage name must not be empty or whitespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateCommonTradeData(tt.tradeData)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateBuyOrSell(t *testing.T) {
	tests := []struct {
		name      string
		tradeInfo *types.TradeInfo
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "valid trade data object",
			tradeInfo: types.GetSampleTradeData(types.TradeTypeSell).TradeInfo,
		},
		{
			name: "invalid share_price",
			tradeInfo: &types.TradeInfo{
				SharePrice: 0,
			},
			expErr:    true,
			expErrMsg: "share_price must be greater than 0",
		},
		{
			name: "invalid share_net_price",
			tradeInfo: &types.TradeInfo{
				SharePrice:    100,
				ShareNetPrice: 0,
			},
			expErr:    true,
			expErrMsg: "share_net_price must be greater than 0",
		},
		{
			name: "invalid number_of_shares",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 0,
			},
			expErr:    true,
			expErrMsg: "number_of_shares must be greater than 0",
		},
		{
			name: "invalid trade value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     -5,
			},
			expErr:    true,
			expErrMsg: "trade_value must be greater than 0",
		},
		{
			name: "invalid trade net value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  0,
			},
			expErr:    true,
			expErrMsg: "trade_net_value must be greater than 0",
		},
		{
			name: "invalid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity:       nil,
			},
			expErr:    true,
			expErrMsg: "invalid quantity",
		},
		{
			name: "invalid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(100),
				},
			},
			expErr:    true,
			expErrMsg: "invalid quantity",
		},
		{
			name: "zero quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(0),
					Denom:  types.DefaultDenom,
				},
			},
			expErr:    true,
			expErrMsg: "zero quantity not allowed",
		},
		{
			name: "invalid denom",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(1000000),
					Denom:  "uggez1",
				},
			},
			expErr:    true,
			expErrMsg: "invalid denom expected: uggz, got: uggez1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateBuyOrSell(tt.tradeInfo)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateReinvestment(t *testing.T) {
	tests := []struct {
		name      string
		tradeInfo *types.TradeInfo
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "valid trade data object",
			tradeInfo: types.GetSampleTradeData(types.TradeTypeReinvestment).TradeInfo,
		},
		{
			name: "invalid share_price",
			tradeInfo: &types.TradeInfo{
				SharePrice: 0,
			},
			expErr:    true,
			expErrMsg: "share_price must be greater than 0",
		},
		{
			name: "invalid share_net_price",
			tradeInfo: &types.TradeInfo{
				SharePrice:    100,
				ShareNetPrice: 0,
			},
			expErr:    true,
			expErrMsg: "share_net_price must be greater than 0",
		},
		{
			name: "invalid number_of_shares",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 0,
			},
			expErr:    true,
			expErrMsg: "number_of_shares must be greater than 0",
		},
		{
			name: "invalid trade value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     -5,
			},
			expErr:    true,
			expErrMsg: "trade_value must be greater than 0",
		},
		{
			name: "invalid trade net value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  0,
			},
			expErr:    true,
			expErrMsg: "trade_net_value must be greater than 0",
		},
		{
			name: "invalid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  5000,
				TradeType:      types.TradeTypeReinvestment,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(1000000),
					Denom:  types.DefaultDenom,
				},
			},
			expErr:    true,
			expErrMsg: "quantity must not be set for trade type TRADE_TYPE_REINVESTMENT",
		},
		{
			name: "set valid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(0),
					Denom:  "",
				},
			},
			expErr: false,
		},
		{
			name: "nil quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     100,
				ShareNetPrice:  100,
				NumberOfShares: 500,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity:       nil,
			},
			expErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateReinvestment(tt.tradeInfo)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateDividends(t *testing.T) {
	tests := []struct {
		name      string
		tradeInfo *types.TradeInfo
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "valid trade data object",
			tradeInfo: types.GetSampleTradeData(types.TradeTypeDividends).TradeInfo,
		},
		{
			name: "invalid share_price",
			tradeInfo: &types.TradeInfo{
				SharePrice: 100,
			},
			expErr:    true,
			expErrMsg: "share_price must be 0",
		},
		{
			name: "invalid share_net_price",
			tradeInfo: &types.TradeInfo{
				SharePrice:    0,
				ShareNetPrice: 100,
			},
			expErr:    true,
			expErrMsg: "share_net_price must be 0",
		},
		{
			name: "invalid number_of_shares",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 100,
			},
			expErr:    true,
			expErrMsg: "number_of_shares must be 0",
		},
		{
			name: "invalid trade value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 0,
				TradeValue:     -5,
			},
			expErr:    true,
			expErrMsg: "trade_value must be greater than 0",
		},
		{
			name: "invalid trade net value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 0,
				TradeValue:     5000,
				TradeNetValue:  0,
			},
			expErr:    true,
			expErrMsg: "trade_net_value must be greater than 0",
		},
		{
			name: "invalid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 0,
				TradeValue:     5000,
				TradeNetValue:  5000,
				TradeType:      types.TradeTypeReinvestment,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(1000000),
					Denom:  types.DefaultDenom,
				},
			},
			expErr:    true,
			expErrMsg: "quantity must not be set for trade type TRADE_TYPE_REINVESTMENT",
		},
		{
			name: "set valid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 0,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(0),
					Denom:  "",
				},
			},
			expErr: false,
		},
		{
			name: "nil quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 0,
				TradeValue:     5000,
				TradeNetValue:  5000,
				Quantity:       nil,
			},
			expErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateDividends(tt.tradeInfo)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSplit(t *testing.T) {
	tests := []struct {
		name      string
		tradeInfo *types.TradeInfo
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "valid trade data object",
			tradeInfo: types.GetSampleTradeData(types.TradeTypeSplit).TradeInfo,
		},
		{
			name: "invalid share_price",
			tradeInfo: &types.TradeInfo{
				SharePrice: -1,
			},
			expErr:    true,
			expErrMsg: "share_price must be 0",
		},
		{
			name: "invalid share_net_price",
			tradeInfo: &types.TradeInfo{
				SharePrice:    0,
				ShareNetPrice: -5,
			},
			expErr:    true,
			expErrMsg: "share_net_price must be 0",
		},
		{
			name: "invalid number_of_shares",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 0,
			},
			expErr:    true,
			expErrMsg: "number_of_shares must be greater than 0",
		},
		{
			name: "invalid trade value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 500,
				TradeValue:     1000,
			},
			expErr:    true,
			expErrMsg: "trade_value must be 0",
		},
		{
			name: "invalid trade net value",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 500,
				TradeValue:     0,
				TradeNetValue:  100,
			},
			expErr:    true,
			expErrMsg: "trade_net_value must be 0",
		},
		{
			name: "invalid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 500,
				TradeValue:     0,
				TradeNetValue:  0,
				TradeType:      types.TradeTypeSplit,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(1000000),
					Denom:  types.DefaultDenom,
				},
			},
			expErr:    true,
			expErrMsg: "quantity must not be set for trade type TRADE_TYPE_SPLIT",
		},
		{
			name: "set valid quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 500,
				TradeValue:     0,
				TradeNetValue:  0,
				Quantity: &sdk.Coin{
					Amount: math.NewInt(0),
					Denom:  "",
				},
			},
			expErr: false,
		},
		{
			name: "nil quantity",
			tradeInfo: &types.TradeInfo{
				SharePrice:     0,
				ShareNetPrice:  0,
				NumberOfShares: 500,
				TradeValue:     0,
				TradeNetValue:  0,
				Quantity:       nil,
			},
			expErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateSplit(tt.tradeInfo)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateNoQuantity(t *testing.T) {
	tests := []struct {
		name      string
		tradeInfo *types.TradeInfo
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid quantity",
			tradeInfo: &types.TradeInfo{
				Quantity: &sdk.Coin{
					Amount: math.NewInt(1000000),
					Denom:  types.DefaultDenom,
				},
			},
			expErr:    true,
			expErrMsg: "quantity must not be set",
		},
		{
			name: "set valid quantity",
			tradeInfo: &types.TradeInfo{
				Quantity: &sdk.Coin{
					Amount: math.NewInt(0),
					Denom:  "",
				},
			},
			expErr: false,
		},
		{
			name: "nil quantity",
			tradeInfo: &types.TradeInfo{
				Quantity: nil,
			},
			expErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.ValidateNoQuantity(tt.tradeInfo)
			if tt.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
