package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgProcessTrade_ValidateBasic(t *testing.T) {
	sdkTypes.GetConfig().SetBech32PrefixForAccount("ggez", "ggez")

	tests := []struct {
		name string
		msg  MsgProcessTrade
		err  error
	}{
		{
			name: "process trade with valid data (confirm)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: Confirm,
				TradeIndex:  1,
			},
			err: nil,
		},
		{
			name: "process trade with valid data (reject)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: Reject,
				TradeIndex:  1,
			},
			err: nil,
		},
		{
			name: "process trade with invalid address",
			msg: MsgProcessTrade{
				Creator:     "xxxx1uuyxga4x50h43yucgtn8ddxd5au5nvh0dlf3fl",
				ProcessType: Confirm,
				TradeIndex:  1,
			},
			err: ErrInvalidCreatorAddress,
		},
		{
			name: "process trade with invalid address (empty)",
			msg: MsgProcessTrade{
				Creator:     "",
				ProcessType: Confirm,
				TradeIndex:  1,
			},
			err: ErrInvalidCreatorAddress,
		},
		{
			name: "process trade with invalid process type",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: "XXXX",
				TradeIndex:  1,
			},
			err: ErrInvalidProcessType,
		},
		{
			name: "process trade with invalid process type (empty)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: "",
				TradeIndex:  1,
			},
			err: ErrInvalidProcessType,
		},
		{
			name: "process trade with invalid trade index (not number)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: Confirm,
				TradeIndex:  0,
			},
			err: ErrInvalidTradeIndex,
		},
		{
			name: "process trade with invalid trade index (negative)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: Reject,
				TradeIndex:  0,
			},
			err: ErrInvalidTradeIndex,
		},
		{
			name: "process trade with invalid trade index (zero)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: Reject,
				TradeIndex:  0,
			},
			err: ErrInvalidTradeIndex,
		},
		{
			name: "process trade with invalid trade index (empty)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: Reject,
				TradeIndex:  0,
			},
			err: ErrInvalidTradeIndex,
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
