package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateSuperAdmin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateSuperAdmin
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgUpdateSuperAdmin{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new super admin address",
			msg: MsgUpdateSuperAdmin{
				Creator:       sample.AccAddress(),
				NewSuperAdmin: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "all good",
			msg: MsgUpdateSuperAdmin{
				Creator:       sample.AccAddress(),
				NewSuperAdmin: sample.AccAddress(),
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
