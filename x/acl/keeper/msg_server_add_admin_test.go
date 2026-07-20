package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddAdmin(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	superAdmin := sample.AccAddress()
	duplicateAdmin := sample.AccAddress()

	require.NoError(t, f.keeper.SuperAdmin.Set(f.ctx, types.SuperAdmin{Admin: superAdmin}))
	require.NoError(t, f.keeper.AclAdmin.Set(f.ctx, duplicateAdmin, types.AclAdmin{Address: duplicateAdmin}))

	testCases := []struct {
		name        string
		input       *types.MsgAddAdmin
		expectedLen int
		expErr      bool
		expErrMsg   string
	}{
		{
			name: "address unauthorized",
			input: &types.MsgAddAdmin{
				Creator: sample.AccAddress(),
				Admins:  []string{sample.AccAddress()},
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "duplicate admin",
			input: &types.MsgAddAdmin{
				Creator: superAdmin,
				Admins:  []string{duplicateAdmin, sample.AccAddress()},
			},
			expErr:    true,
			expErrMsg: "admin(s) already exist",
		},
		{
			name: "all good",
			input: &types.MsgAddAdmin{
				Creator: superAdmin,
				Admins:  []string{sample.AccAddress(), sample.AccAddress()},
			},
			// duplicateAdmin + 2
			expectedLen: 3,
			expErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.AddAdmin(f.ctx, tc.input)

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
