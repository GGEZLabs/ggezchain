package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestMsgInit(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	testCases := []struct {
		name      string
		input     *types.MsgInit
		expErr    bool
		expErrMsg string
	}{
		{
			name: "all good",
			input: &types.MsgInit{
				Creator:    sample.AccAddress(),
				SuperAdmin: sample.AccAddress(),
			},
			expErr: false,
		},
		{
			name: "super admin initialized",
			input: &types.MsgInit{
				Creator:    sample.AccAddress(),
				SuperAdmin: sample.AccAddress(),
			},
			expErr:    true,
			expErrMsg: "super admin already initialized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.Init(f.ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
