package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddAdmin_ValidateBasic(t *testing.T) {
	duplicateAdmin := sample.AccAddress()

	tests := []struct {
		name string
		msg  MsgAddAdmin
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgAddAdmin{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "zero admins",
			msg: MsgAddAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid admin address",
			msg: MsgAddAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{"invalid_address", sample.AccAddress()},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "duplicate admin address",
			msg: MsgAddAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{duplicateAdmin, duplicateAdmin},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "all good",
			msg: MsgAddAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{sample.AccAddress(), sample.AccAddress()},
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
