package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDeleteAdmin(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	superAdmin := sample.AccAddress()
	alice := sample.AccAddress()
	bob := sample.AccAddress()

	require.NoError(t, f.keeper.SuperAdmin.Set(f.ctx, types.SuperAdmin{Admin: superAdmin}))
	aclAdmins := types.ConvertStringsToAclAdmins([]string{alice, bob})
	for _, aclAdmin := range aclAdmins {
		require.NoError(t, f.keeper.AclAdmin.Set(f.ctx, aclAdmin.Address, aclAdmin))
	}

	testCases := []struct {
		name        string
		input       *types.MsgDeleteAdmin
		expectedLen int
		expErr      bool
		expErrMsg   string
	}{
		{
			name: "address unauthorized",
			input: &types.MsgDeleteAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{alice},
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "admin not exist",
			input: &types.MsgDeleteAdmin{
				Creator: superAdmin,
				Admins:  []string{sample.AccAddress()},
			},
			expErr:    true,
			expErrMsg: "admin(s) does not exist",
		},
		{
			name: "all good",
			input: &types.MsgDeleteAdmin{
				Creator: superAdmin,
				Admins:  []string{alice},
			},
			expectedLen: 1,
			expErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.DeleteAdmin(f.ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				all, err := f.keeper.GetAllAclAdmin(f.ctx)
				require.NoError(t, err)
				require.Len(t, all, tc.expectedLen)
			}
		})
	}
}
