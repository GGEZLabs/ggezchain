package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDeleteAuthority(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	superAdmin := sample.AccAddress()
	admin := sample.AccAddress()
	alice := sample.AccAddress()
	bob := sample.AccAddress()
	aclAuthorityAlice := types.AclAuthority{
		Address:           alice,
		Name:              "Alice",
		AccessDefinitions: []*types.AccessDefinition{},
	}
	aclAuthorityBob := types.AclAuthority{
		Address:           bob,
		Name:              "Bob",
		AccessDefinitions: []*types.AccessDefinition{},
	}
	require.NoError(t, f.keeper.SuperAdmin.Set(f.ctx, types.SuperAdmin{Admin: superAdmin}))
	require.NoError(t, f.keeper.AclAuthority.Set(f.ctx, aclAuthorityAlice.Address, aclAuthorityAlice))
	require.NoError(t, f.keeper.AclAuthority.Set(f.ctx, aclAuthorityBob.Address, aclAuthorityBob))
	require.NoError(t, f.keeper.AclAdmin.Set(f.ctx, admin, types.AclAdmin{Address: admin}))

	testCases := []struct {
		name      string
		input     *types.MsgDeleteAuthority
		expErr    bool
		expErrMsg string
	}{
		{
			name: "unauthorized account",
			input: &types.MsgDeleteAuthority{
				Creator:     sample.AccAddress(),
				AuthAddress: alice,
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "authority not found",
			input: &types.MsgDeleteAuthority{
				Creator:     admin,
				AuthAddress: sample.AccAddress(),
			},
			expErr:    true,
			expErrMsg: "authority address does not exist",
		},
		{
			name: "delete authority by super admin",
			input: &types.MsgDeleteAuthority{
				Creator:     superAdmin,
				AuthAddress: alice,
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name: "delete authority by admin",
			input: &types.MsgDeleteAuthority{
				Creator:     admin,
				AuthAddress: bob,
			},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.DeleteAuthority(f.ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
