package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDeleteAuthority(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	admin := sample.AccAddress()
	alice := sample.AccAddress()
	bob := sample.AccAddress()
	aclAuthority := types.AclAuthority{
		Address:           alice,
		Name:              "Alice",
		AccessDefinitions: []*types.AccessDefinition{},
	}
	k.SetAclAuthority(ctx, aclAuthority)
	k.SetAclAdmin(ctx, types.AclAdmin{Address: admin})
	wctx := sdk.UnwrapSDKContext(ctx)

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
				AuthAddress: bob,
			},
			expErr:    true,
			expErrMsg: "authority address does not exist",
		},
		{
			name: "all good",
			input: &types.MsgDeleteAuthority{
				Creator:     admin,
				AuthAddress: alice,
			},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.DeleteAuthority(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
