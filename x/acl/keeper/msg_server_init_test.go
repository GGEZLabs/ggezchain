package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgInit(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	aclAdmin := sample.AccAddress()
	wctx := sdk.UnwrapSDKContext(ctx)

	testCases := []struct {
		name      string
		input     *types.MsgInit
		fun       func()
		expErr    bool
		expErrMsg string
	}{
		{
			name: "acl admin initialized",
			input: &types.MsgInit{
				Creator: sample.AccAddress(),
				Admins:  []string{sample.AccAddress()},
			},
			fun: func() {
				k.SetAclAdmin(ctx, types.AclAdmin{Address: aclAdmin})
			},
			expErr:    true,
			expErrMsg: "acl admin already initialized",
		},
		{
			name: "all good",
			input: &types.MsgInit{
				Creator: sample.AccAddress(),
				Admins:  []string{sample.AccAddress()},
			},
			fun: func() {
				k.RemoveAclAdmin(ctx, aclAdmin)
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fun()
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
