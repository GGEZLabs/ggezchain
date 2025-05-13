package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgProcessTrade_ValidateBasic(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("ggez", "ggez")

	tests := []struct {
		name string
		msg  MsgProcessTrade
		err  error
	}{
		{
			name: "process trade with valid data (confirm)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: ProcessTypeConfirm,
				TradeIndex:  1,
			},
		},
		{
			name: "process trade with valid data (reject)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: ProcessTypeReject,
				TradeIndex:  1,
			},
		},
		{
			name: "process trade with invalid address",
			msg: MsgProcessTrade{
				Creator:     "xxxx1uuyxga4x50h43yucgtn8ddxd5au5nvh0dlf3fl",
				ProcessType: ProcessTypeConfirm,
				TradeIndex:  1,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "process trade with invalid address (empty)",
			msg: MsgProcessTrade{
				Creator:     "",
				ProcessType: ProcessTypeConfirm,
				TradeIndex:  1,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "process trade with invalid process type",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: ProcessTypeUnspecified,
				TradeIndex:  1,
			},
			err: ErrInvalidProcessType,
		},
		{
			name: "process trade with invalid trade index (zero)",
			msg: MsgProcessTrade{
				Creator:     sample.AccAddress(),
				ProcessType: ProcessTypeReject,
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
