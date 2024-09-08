package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/stretchr/testify/require"
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
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: nil,
		},
		{
			name: "create trade with valid data (sell)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Sell,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: nil,
		},
		{
			name: "create trade with invalid address",
			msg: MsgCreateTrade{
				Creator:              "xxxx1uuyxga4x50h43yucgtn8ddxd5au5nvh0dlf3fl",
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidCreatorAddress,
		},
		{
			name: "create trade with invalid address (empty)",
			msg: MsgCreateTrade{
				Creator:              "",
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidCreatorAddress,
		},
		{
			name: "create trade with invalid trade type",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            "XXXX",
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradeType,
		},
		{
			name: "create trade with invalid trade type (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            "",
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradeType,
		},
		{
			name: "create trade with invalid coin denom",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 "XXXX",
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidCoinDenom,
		},
		{
			name: "create trade with invalid coin denom (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 "",
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidCoinDenom,
		},
		{
			name: "create trade with invalid trade price (not number)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "XXXX",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid trade price (negative)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "-0.001",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid trade price (zero)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid trade price (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradePrice,
		},
		{
			name: "create trade with invalid quantity (not number)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "XXXX",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid quantity (negative)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "-100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid quantity (zero)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "0",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid quantity (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradeQuantity,
		},
		{
			name: "create trade with invalid receiver address",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      "xxxx1uuyxga4x50h43yucgtn8ddxd5au5nvh0dlf3fl",
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidReceiverAddress,
		},
		{
			name: "create trade with invalid receiver address (empty)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "0.001",
				Quantity:             "100000",
				ReceiverAddress:      "",
				TradeData:            "TradeData",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidReceiverAddress,
		},
		{
			name: "create trade with invalid trade data",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "10",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "",
				BankingSystemData:    "{}",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
			},
			err: ErrInvalidTradeData,
		},
		{
			name: "create trade with invalid banking system data",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				TradeType:            Buy,
				Coin:                 DefaultCoinDenom,
				Price:                "10",
				Quantity:             "100000",
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            "TradeData",
				BankingSystemData:    "",
				CoinMintingPriceJSON: "{}",
				ExchangeRateJSON:     "{}",
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
