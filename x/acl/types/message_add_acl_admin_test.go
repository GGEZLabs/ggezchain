package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgAddAclAdmin_ValidateBasic(t *testing.T) {
	duplicateAdmin := sample.AccAddress()

	tests := []struct {
		name string
		msg  MsgAddAclAdmin
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgAddAclAdmin{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "zero admins",
			msg: MsgAddAclAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid admin address",
			msg: MsgAddAclAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{"invalid_address", sample.AccAddress()},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "duplicate admin address",
			msg: MsgAddAclAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{duplicateAdmin, duplicateAdmin},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "all good",
			msg: MsgAddAclAdmin{
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
