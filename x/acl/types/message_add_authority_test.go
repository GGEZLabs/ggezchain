package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAddAuthority_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddAuthority
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgAddAuthority{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid auth address",
			msg: MsgAddAuthority{
				Creator:     sample.AccAddress(),
				AuthAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty name",
			msg: MsgAddAuthority{
				Creator:     sample.AccAddress(),
				AuthAddress: sample.AccAddress(),
				Name:        "",
			},
			err: ErrEmptyName,
		},
		{
			name: "invalid access definitions",
			msg: MsgAddAuthority{
				Creator:           sample.AccAddress(),
				AuthAddress:       sample.AccAddress(),
				Name:              "Alice",
				AccessDefinitions: `{"module":"module1","is_maker":true "is_checker":true}`,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "all good",
			msg: MsgAddAuthority{
				Creator:           sample.AccAddress(),
				AuthAddress:       sample.AccAddress(),
				Name:              "Alice",
				AccessDefinitions: `[{"module":"module1","is_maker":true, "is_checker":true}]`,
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
