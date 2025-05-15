package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgDeleteAclAdmin(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	aclAdmin := sample.AccAddress()
	alice := sample.AccAddress()
	bob := sample.AccAddress()

	aclAdmins := types.ConvertStringsToAclAdmins([]string{aclAdmin, alice, bob})
	k.SetAclAdmins(ctx, aclAdmins)

	testCases := []struct {
		name        string
		input       *types.MsgDeleteAclAdmin
		expectedLen int
		expErr      bool
		expErrMsg   string
	}{
		{
			name: "address unauthorized",
			input: &types.MsgDeleteAclAdmin{
				Creator: sample.AccAddress(),
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "delete all admins",
			input: &types.MsgDeleteAclAdmin{
				Creator: aclAdmin,
				Admins:  []string{aclAdmin, alice, bob},
			},
			expErr:    true,
			expErrMsg: "cannot delete all admins, at least one aclAdmin must remain",
		},
		{
			name: "admin not exist",
			input: &types.MsgDeleteAclAdmin{
				Creator: aclAdmin,
				Admins:  []string{sample.AccAddress()},
			},
			expErr:    true,
			expErrMsg: "admin(s) not exist",
		},
		{
			name: "all good",
			input: &types.MsgDeleteAclAdmin{
				Creator: aclAdmin,
				Admins:  []string{alice, bob},
			},
			expectedLen: 1,
			expErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.DeleteAclAdmin(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.Equal(t, len(k.GetAllAclAdmin(ctx)), tc.expectedLen)
				require.NoError(t, err)
			}
		})
	}
}
