package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name      string
		msg       MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority address",
			msg: MsgUpdateParams{
				Authority: "invalid_address",
			},
			expErr:    true,
			expErrMsg: "invalid authority address",
		},
		{
			name: "all good",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params:    DefaultParams(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.expErr {
				require.Contains(t, err.Error(), tt.expErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
