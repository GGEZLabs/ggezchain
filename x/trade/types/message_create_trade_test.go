package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateTrade_ValidateBasic(t *testing.T) {
	td := GetSampleTradeData(TradeTypeBuy)
	tests := []struct {
		name string
		msg  MsgCreateTrade
		err  error
	}{
		{
			name: "create trade with valid data (buy)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            td,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
		},
		{
			name: "create trade with valid data (sell)",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            GetSampleTradeData(TradeTypeSell),
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
			},
		},
		{
			name: "create trade with invalid creator address",
			msg: MsgCreateTrade{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "create trade with invalid address (empty)",
			msg: MsgCreateTrade{
				Creator: "",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "create trade with invalid trade data",
			msg: MsgCreateTrade{
				Creator:         sample.AccAddress(),
				ReceiverAddress: sample.AccAddress(),
				TradeData:       `"trade_data":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"currency":"USD","exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"price":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}`,
			},
			err: ErrInvalidTradeData,
		},
		{
			name: "create trade with invalid banking system data",
			msg: MsgCreateTrade{
				Creator:           sample.AccAddress(),
				ReceiverAddress:   sample.AccAddress(),
				TradeData:         td,
				BankingSystemData: "",
			},
			err: ErrInvalidBankingSystemData,
		},
		{
			name: "create trade with valid create date",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            td,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
				CreateDate:           "2023-05-11T08:44:00Z",
			},
		},
		{
			name: "create trade with invalid create date",
			msg: MsgCreateTrade{
				Creator:              sample.AccAddress(),
				ReceiverAddress:      sample.AccAddress(),
				TradeData:            td,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: "{}",
				ExchangeRateJson:     "{}",
				CreateDate:           "2023-05-11",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
