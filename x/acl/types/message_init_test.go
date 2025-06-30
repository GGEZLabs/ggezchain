package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgInit_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgInit
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgInit{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid super admin address",
			msg: MsgInit{
				Creator:    sample.AccAddress(),
				SuperAdmin: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "all good",
			msg: MsgInit{
				Creator:    sample.AccAddress(),
				SuperAdmin: sample.AccAddress(),
			},
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
