package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgUpdateParams(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))
	wctx := sdk.UnwrapSDKContext(ctx)

	// default params
	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgUpdateParams{
				Authority: "invalid",
				Params:    params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		// {
		// 	name: "invalid admin",
		// 	input: &types.MsgUpdateParams{
		// 		Authority: k.GetAuthority(),
		// 		Params:    types.Params{
		// 			Admin: "invalid_address",
		// 		},
		// 	},
		// 	expErr:    true,
		// 	expErrMsg: "invalid admin address",
		// },
		// {
		// 	name: "empty admin",
		// 	input: &types.MsgUpdateParams{
		// 		Authority: k.GetAuthority(),
		// 		Params:    types.Params{
		// 			Admin: "",
		// 		},
		// 	},
		// 	expErr:    true,
		// 	expErrMsg: "admin address cannot be empty",
		// },
		{
			name: "all good",
			input: &types.MsgUpdateParams{
				Authority: k.GetAuthority(),
				Params:    params,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdateParams(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
