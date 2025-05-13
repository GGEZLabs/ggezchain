package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgCreateTrade_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateTrade
		err  error
	}{
		{
			name: "create trade with valid data (buy)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
		},
		{
			name: "create trade with valid data (sell)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeSell,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
		},
		{
			name: "create trade with invalid creator address",
			msg: MsgCreateTrade{
				Creator:              "invalid_address",
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "create trade with invalid address (empty)",
			msg: MsgCreateTrade{
				Creator:              "",
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "create trade with invalid trade type",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeUnspecified,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradeType,
		},
		{
			name: "create trade with invalid coin denom",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: "invalid_denom", Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidDenom,
		},
		{
			name: "create trade with invalid coin denom (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: "", Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidDenom,
		},
		{
			name: "create trade with invalid trade price (not number)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "XXXX",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid trade price (negative)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "-0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid trade price (zero)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid trade price (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid quantity (not number)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid quantity (negative)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid quantity (zero)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid quantity (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid receiver address",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      "invalid_address",
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "create trade with invalid receiver address (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "0.001",
				ReceiverAddress:      "",
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "create trade with invalid trade data",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "10",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidTradeData,
		},
		{
			name: "create trade with invalid banking system data",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            TradeTypeBuy,
				Amount:               &types.Coin{Denom: DefaultCoinDenom, Amount: math.NewInt(100000)},
				Price:                "10",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            `{"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":"buy","trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":162075000000000,"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
				BankingSystemData:    "",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
			err: ErrInvalidBankingSystemData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
