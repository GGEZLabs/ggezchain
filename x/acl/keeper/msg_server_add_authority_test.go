package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddAuthority(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	superAdmin := sample.AccAddress()
	admin := sample.AccAddress()
	alice := sample.AccAddress()
	bob := sample.AccAddress()
	sampleAccessDefinitions := `[{"module":"trade","is_maker":true,"is_checker":true}]`
	aclAuthority := types.AclAuthority{
		Address:           alice,
		Name:              "Alice",
		AccessDefinitions: []*types.AccessDefinition{},
	}
	k.SetSuperAdmin(ctx, types.SuperAdmin{Admin: superAdmin})
	k.SetAclAuthority(ctx, aclAuthority)
	k.SetAclAdmin(ctx, types.AclAdmin{Address: admin})
	wctx := sdk.UnwrapSDKContext(ctx)

	testCases := []struct {
		name      string
		input     *types.MsgAddAuthority
		expErr    bool
		expErrMsg string
	}{
		{
			name: "unauthorized account",
			input: &types.MsgAddAuthority{
				Creator:           sample.AccAddress(),
				AuthAddress:       alice,
				Name:              "Alice",
				AccessDefinitions: sampleAccessDefinitions,
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "authority already exist",
			input: &types.MsgAddAuthority{
				Creator:           admin,
				AuthAddress:       alice,
				Name:              "Alice",
				AccessDefinitions: sampleAccessDefinitions,
			},
			expErr:    true,
			expErrMsg: "authority address already exist",
		},
		{
			name: "invalid access definitions format",
			input: &types.MsgAddAuthority{
				Creator:           admin,
				AuthAddress:       bob,
				Name:              "Bob",
				AccessDefinitions: `[{"module":"trade","is_maker":true "is_checker":true}]`,
			},
			expErr:    true,
			expErrMsg: "invalid access definition list format",
		},
		{
			name: "empty access definition list",
			input: &types.MsgAddAuthority{
				Creator:           admin,
				AuthAddress:       bob,
				Name:              "Bob",
				AccessDefinitions: `[]`,
			},
			expErr:    true,
			expErrMsg: "access definition list is required and cannot be empty",
		},
		{
			name: "add empty module",
			input: &types.MsgAddAuthority{
				Creator:           admin,
				AuthAddress:       bob,
				Name:              "Bob",
				AccessDefinitions: `[{"module":"","is_maker":true,"is_checker":true}]`,
			},
			expErr:    true,
			expErrMsg: "invalid module name",
		},
		{
			name: "add duplicated modules",
			input: &types.MsgAddAuthority{
				Creator:           admin,
				AuthAddress:       bob,
				Name:              "Bob",
				AccessDefinitions: `[{"module":"trade","is_maker":true,"is_checker":true},{"module":"trade","is_maker":true,"is_checker":true}]`,
			},
			expErr:    true,
			expErrMsg: "invalid module name",
		},
		{
			name: "add authority by super admin",
			input: &types.MsgAddAuthority{
				Creator:           superAdmin,
				AuthAddress:       bob,
				Name:              "Bob",
				AccessDefinitions: sampleAccessDefinitions,
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name: "add authority by admin",
			input: &types.MsgAddAuthority{
				Creator:           admin,
				AuthAddress:       sample.AccAddress(),
				Name:              "Carol",
				AccessDefinitions: sampleAccessDefinitions,
			},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.AddAuthority(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
