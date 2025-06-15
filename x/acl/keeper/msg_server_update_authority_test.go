package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateAuthority(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	admin := sample.AccAddress()
	alice := sample.AccAddress()
	bob := sample.AccAddress()
	aclAuthority := types.AclAuthority{
		Address: alice,
		Name:    "alice",
		AccessDefinitions: []*types.AccessDefinition{
			{Module: "module1", IsMaker: true, IsChecker: false},
			{Module: "module2", IsMaker: true, IsChecker: false},
		},
	}
	k.SetAclAuthority(ctx, aclAuthority)
	k.SetAclAdmin(ctx, types.AclAdmin{Address: admin})
	wctx := sdk.UnwrapSDKContext(ctx)

	testCases := []struct {
		name      string
		input     *types.MsgUpdateAuthority
		expErr    bool
		expErrMsg string
	}{
		{
			name: "unauthorized account",
			input: &types.MsgUpdateAuthority{
				Creator:     sample.AccAddress(),
				AuthAddress: alice,
				NewName:     "Carol",
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
		{
			name: "authority not found",
			input: &types.MsgUpdateAuthority{
				Creator:     admin,
				AuthAddress: bob,
			},
			expErr:    true,
			expErrMsg: "authority address not exist",
		},
		{
			name: "empty access definition list",
			input: &types.MsgUpdateAuthority{
				Creator:                    admin,
				AuthAddress:                alice,
				OverwriteAccessDefinitions: `[]`,
			},
			expErr:    true,
			expErrMsg: "access definition list is empty",
		},
		{
			name: "invalid overwrite access definitions format",
			input: &types.MsgUpdateAuthority{
				Creator:                    admin,
				AuthAddress:                alice,
				OverwriteAccessDefinitions: `[{"module":"trade","is_maker":true "is_checker":true}]`,
			},
			expErr:    true,
			expErrMsg: "invalid access definition list format",
		},
		{
			name: "invalid overwrite access definitions (empty module)",
			input: &types.MsgUpdateAuthority{
				Creator:                    admin,
				AuthAddress:                alice,
				OverwriteAccessDefinitions: `[{"module":"","is_maker":false,"is_checker":true}]`,
			},
			expErr:    true,
			expErrMsg: "invalid module name",
		},
		{
			name: "invalid overwrite access definitions (duplicate module)",
			input: &types.MsgUpdateAuthority{
				Creator:                    admin,
				AuthAddress:                alice,
				OverwriteAccessDefinitions: `[{"module":"module1","is_maker":false,"is_checker":true},{"module":"module1","is_maker":false,"is_checker":true}]`,
			},
			expErr:    true,
			expErrMsg: "module1 module(s) is duplicates",
		},
		{
			name: "invalid update access definitions format",
			input: &types.MsgUpdateAuthority{
				Creator:                admin,
				AuthAddress:            alice,
				UpdateAccessDefinition: `{"module":"trade","is_maker":true "is_checker":true}`,
			},
			expErr:    true,
			expErrMsg: "invalid access definition object format",
		},
		{
			name: "invalid add access definitions format",
			input: &types.MsgUpdateAuthority{
				Creator:              admin,
				AuthAddress:          alice,
				AddAccessDefinitions: `[{"module":"trade","is_maker":true "is_checker":true}]`,
			},
			expErr:    true,
			expErrMsg: "invalid access definition list format",
		},
		{
			name: "invalid delete access definitions (module not exist)",
			input: &types.MsgUpdateAuthority{
				Creator:                 admin,
				AuthAddress:             alice,
				DeleteAccessDefinitions: []string{"trade"},
			},
			expErr:    true,
			expErrMsg: "module not exist",
		},
		{
			name: "invalid delete access definitions (empty module)",
			input: &types.MsgUpdateAuthority{
				Creator:                 admin,
				AuthAddress:             alice,
				DeleteAccessDefinitions: []string{""},
			},
			expErr:    true,
			expErrMsg: "module name cannot be empty",
		},
		{
			name: "invalid delete access definitions (duplicate module)",
			input: &types.MsgUpdateAuthority{
				Creator:                 admin,
				AuthAddress:             alice,
				DeleteAccessDefinitions: []string{"module1", "module1"},
			},
			expErr:    true,
			expErrMsg: "module1 module(s) is duplicates",
		},
		{
			name: "update and delete same module",
			input: &types.MsgUpdateAuthority{
				Creator:                 admin,
				AuthAddress:             alice,
				UpdateAccessDefinition:  `{"module":"module1","is_maker":true, "is_checker":true}`,
				DeleteAccessDefinitions: []string{"module1"},
			},
			expErr:    true,
			expErrMsg: "module(s) cannot be both added/updated and removed in the same request",
		},
		{
			name: "add and delete same module",
			input: &types.MsgUpdateAuthority{
				Creator:                 admin,
				AuthAddress:             alice,
				AddAccessDefinitions:    `[{"module":"module1","is_maker":true, "is_checker":true}]`,
				DeleteAccessDefinitions: []string{"module1"},
			},
			expErr:    true,
			expErrMsg: "module(s) cannot be both added/updated and removed in the same request",
		},
		{
			name: "update add and delete same module",
			input: &types.MsgUpdateAuthority{
				Creator:                 admin,
				AuthAddress:             alice,
				AddAccessDefinitions:    `[{"module":"module1","is_maker":true, "is_checker":true}]`,
				UpdateAccessDefinition:  `{"module":"module1","is_maker":true, "is_checker":true}`,
				DeleteAccessDefinitions: []string{"module1"},
			},
			expErr:    true,
			expErrMsg: "module(s) cannot be both added/updated and removed in the same request",
		},
		{
			name: "update name only",
			input: &types.MsgUpdateAuthority{
				Creator:     admin,
				AuthAddress: alice,
				NewName:     "Carol",
			},
			expErr: false,
		},
		{
			name: "add, update and delete definitions",
			input: &types.MsgUpdateAuthority{
				Creator:                 admin,
				AuthAddress:             alice,
				NewName:                 "new name",
				UpdateAccessDefinition:  `{"module":"module2","is_maker":false,"is_checker":true}`,
				AddAccessDefinitions:    `[{"module":"module3","is_maker":true,"is_checker":true}]`,
				DeleteAccessDefinitions: []string{"module1"},
			},
			expErr: false,
		},
		{
			name: "overwrite access definitions",
			input: &types.MsgUpdateAuthority{
				Creator:                    admin,
				AuthAddress:                alice,
				OverwriteAccessDefinitions: `[{"module":"module1","is_maker":true,"is_checker":false}]`,
			},
			expErr: false,
		},
		{
			name: "clear all access definitions",
			input: &types.MsgUpdateAuthority{
				Creator:                   admin,
				AuthAddress:               alice,
				ClearAllAccessDefinitions: true,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdateAuthority(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
