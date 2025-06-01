package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgInit(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	testCases := []struct {
		name      string
		input     *types.MsgInit
		expErr    bool
		expErrMsg string
	}{
		{
			name: "all good",
			input: &types.MsgInit{
				Creator:    sample.AccAddress(),
				SuperAdmin: sample.AccAddress(),
			},
			expErr: false,
		},
		{
			name: "super admin initialized",
			input: &types.MsgInit{
				Creator:    sample.AccAddress(),
				SuperAdmin: sample.AccAddress(),
			},
			expErr:    true,
			expErrMsg: "super admin already initialized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.Init(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
