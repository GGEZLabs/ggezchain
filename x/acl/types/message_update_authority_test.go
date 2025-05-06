package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateAuthority_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateAuthority
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgUpdateAuthority{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid auth address",
			msg: MsgUpdateAuthority{
				Creator:     sample.AccAddress(),
				AuthAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "no updates",
			msg: MsgUpdateAuthority{
				Creator:     sample.AccAddress(),
				AuthAddress: sample.AccAddress(),
			},
			err: ErrNoUpdateFlags,
		},
		{
			name: "use OverwriteAccessDefinitions with another flag",
			msg: MsgUpdateAuthority{
				Creator:                    sample.AccAddress(),
				AuthAddress:                sample.AccAddress(),
				OverwriteAccessDefinitions: `[{"module":"module1","is_maker":true, "is_checker":true}]`,
				ClearAllAccessDefinitions:  true,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid OverwriteAccessDefinitions JSON format",
			msg: MsgUpdateAuthority{
				Creator:                    sample.AccAddress(),
				AuthAddress:                sample.AccAddress(),
				OverwriteAccessDefinitions: `[{"module":"module1","is_maker":true "is_checker":true}]`,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "use ClearAllAccessDefinitions with another flag",
			msg: MsgUpdateAuthority{
				Creator:                   sample.AccAddress(),
				AuthAddress:               sample.AccAddress(),
				ClearAllAccessDefinitions: true,
				AddAccessDefinitions:      `[{"module":"module1","is_maker":true, "is_checker":true}]`,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid UpdateAccessDefinition JSON format",
			msg: MsgUpdateAuthority{
				Creator:                sample.AccAddress(),
				AuthAddress:            sample.AccAddress(),
				UpdateAccessDefinition: `{"module":"module1","is_maker":true "is_checker":true}`,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid AddAccessDefinitions JSON format",
			msg: MsgUpdateAuthority{
				Creator:              sample.AccAddress(),
				AuthAddress:          sample.AccAddress(),
				AddAccessDefinitions: `[{"module":"module1","is_maker":true "is_checker":true}]`,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "overwrite access definitions",
			msg: MsgUpdateAuthority{
				Creator:                    sample.AccAddress(),
				AuthAddress:                sample.AccAddress(),
				OverwriteAccessDefinitions: `[{"module":"module1","is_maker":true, "is_checker":true}]`,
			},
		},
		{
			name: "clear access definitions",
			msg: MsgUpdateAuthority{
				Creator:                   sample.AccAddress(),
				AuthAddress:               sample.AccAddress(),
				ClearAllAccessDefinitions: true,
			},
		},
		{
			name: "update access definitions",
			msg: MsgUpdateAuthority{
				Creator:                sample.AccAddress(),
				AuthAddress:            sample.AccAddress(),
				UpdateAccessDefinition: `{"module":"module1","is_maker":true, "is_checker":true}`,
			},
		},
		{
			name: "add access definitions",
			msg: MsgUpdateAuthority{
				Creator:              sample.AccAddress(),
				AuthAddress:          sample.AccAddress(),
				AddAccessDefinitions: `[{"module":"module1","is_maker":true, "is_checker":true}]`,
			},
		},
		{
			name: "delete access definitions",
			msg: MsgUpdateAuthority{
				Creator:                 sample.AccAddress(),
				AuthAddress:             sample.AccAddress(),
				DeleteAccessDefinitions: []string{"module1"},
			},
		},
		{
			name: "add update and delete access definitions",
			msg: MsgUpdateAuthority{
				Creator:                 sample.AccAddress(),
				AuthAddress:             sample.AccAddress(),
				AddAccessDefinitions:    `[{"module":"module1","is_maker":true, "is_checker":true}]`,
				UpdateAccessDefinition:  `{"module":"module2","is_maker":true, "is_checker":true}`,
				DeleteAccessDefinitions: []string{"module3"},
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
