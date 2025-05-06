package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("ggez", "pub")
	tests := []struct {
		name          string
		msg           MsgUpdateParams
		expectedError bool
		errMsg        string
	}{
		{
			name: "invalid authority address",
			msg: MsgUpdateParams{
				Authority: "invalid_address",
			},
			expectedError: true,
			errMsg:        "invalid authority address",
		},
		{
			name: "empty admin address",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params:    Params{Admin: ""},
			},
			expectedError: true,
			errMsg:        "admin address cannot be empty",
		},
		{
			name: "invalid admin address",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params:    Params{Admin: "invalid_address"},
			},
			expectedError: true,
			errMsg:        "invalid admin address",
		},
		{
			name: "invalid admin address prefix",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params:    Params{Admin: "cosmos1q3sfaepes35ly4sa5ppguf6gs49un4uzxrupy2"},
			},
			expectedError: true,
			errMsg:        "invalid admin address",
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
			if tt.expectedError {
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
