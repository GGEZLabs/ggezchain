package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgInitAclAdmin(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	aclAdmin := sample.AccAddress()
	wctx := sdk.UnwrapSDKContext(ctx)

	testCases := []struct {
		name      string
		input     *types.MsgInitAclAdmin
		fun       func()
		expErr    bool
		expErrMsg string
	}{
		{
			name: "acl admin initialized",
			input: &types.MsgInitAclAdmin{
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
			input: &types.MsgInitAclAdmin{
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
			_, err := ms.InitAclAdmin(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
