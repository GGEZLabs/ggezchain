package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateSuperAdmin(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)
	superAdmin := sample.AccAddress()
	testCases := []struct {
		name      string
		input     *types.MsgUpdateSuperAdmin
		fun       func()
		expErr    bool
		expErrMsg string
	}{
		{
			name: "super admin not initialized",
			input: &types.MsgUpdateSuperAdmin{
				Creator:       sample.AccAddress(),
				NewSuperAdmin: sample.AccAddress(),
			},
			fun:       func() {},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "unauthorized account",
			input: &types.MsgUpdateSuperAdmin{
				Creator:       sample.AccAddress(),
				NewSuperAdmin: sample.AccAddress(),
			},
			fun: func() {
				require.NoError(t, f.keeper.SuperAdmin.Set(f.ctx, types.SuperAdmin{Admin: superAdmin}))
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "all good",
			input: &types.MsgUpdateSuperAdmin{
				Creator:       superAdmin,
				NewSuperAdmin: sample.AccAddress(),
			},
			fun:    func() {},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fun()
			_, err := ms.UpdateSuperAdmin(f.ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
