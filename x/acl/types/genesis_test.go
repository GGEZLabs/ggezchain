package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("ggez", "pub")
	duplicateAddress := sample.AccAddress()
	tests := []struct {
		desc          string
		genState      *types.GenesisState
		expectedError bool
		errMsg        string
	}{
		{
			desc: "duplicated aclAuthority",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: duplicateAddress,
					},
					{
						Address: duplicateAddress,
					},
				},
				Params: types.Params{Admin: sample.AccAddress()},
			},
			expectedError: false,
			errMsg:        "duplicated index for aclAuthority",
		},
		{
			desc: "invalid aclAuthority address",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: "invalid_address",
					},
					{
						Address: duplicateAddress,
					},
				},
				Params: types.Params{Admin: sample.AccAddress()},
			},
			expectedError: false,
			errMsg:        "invalid address for aclAuthority",
		},
		{
			desc: "invalid admin address",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: sample.AccAddress(),
					},
					{
						Address: sample.AccAddress(),
					},
				},
				Params: types.Params{Admin: "invalid_address"},
			},
			expectedError: false,
			errMsg:        "invalid admin address",
		},
		{
			desc: "empty admin address",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: sample.AccAddress(),
					},
					{
						Address: sample.AccAddress(),
					},
				},
				Params: types.Params{Admin: ""},
			},
			expectedError: false,
			errMsg:        "admin address cannot be empty",
		},
		{
			desc:          "default is valid",
			genState:      types.DefaultGenesis(),
			expectedError: true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: sample.AccAddress(),
						Name:    "Alice",
						AccessDefinitions: []*types.AccessDefinition{
							{Module: "module1", IsMaker: true, IsChecker: false},
							{Module: "module2", IsMaker: true, IsChecker: false},
							{Module: "module3", IsMaker: true, IsChecker: false},
							{Module: "module4", IsMaker: true, IsChecker: false},
						},
					},
					{
						Address: sample.AccAddress(),
						Name:    "Bob",
						AccessDefinitions: []*types.AccessDefinition{
							{Module: "module1", IsMaker: true, IsChecker: false},
							{Module: "module2", IsMaker: true, IsChecker: false},
						},
					},
				},
				Params: types.Params{Admin: sample.AccAddress()},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			expectedError: true,
		},

		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedError {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			}
		})
	}
}
