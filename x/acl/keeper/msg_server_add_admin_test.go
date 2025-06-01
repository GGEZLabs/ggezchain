package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddAdmin(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	superAdmin := sample.AccAddress()
	duplicateAdmin := sample.AccAddress()

	k.SetSuperAdmin(ctx, types.SuperAdmin{Admin: superAdmin})
	k.SetAclAdmin(ctx, types.AclAdmin{Address: duplicateAdmin})

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
			_, err := ms.AddAdmin(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.Len(t, k.GetAllAclAdmin(ctx), tc.expectedLen)
				require.NoError(t, err)
			}
		})
	}
}
