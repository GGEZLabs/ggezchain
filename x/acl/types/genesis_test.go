package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	duplicateAddress := sample.AccAddress()
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
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
			},
			expErr:    true,
			expErrMsg: "duplicated index for aclAuthority",
		},
		{
			desc: "invalid aclAuthority address",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: "invalid_address",
					},
					{
						Address: sample.AccAddress(),
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid address for aclAuthority",
		},
		{
			desc: "duplicated aclAdmin",
			genState: &types.GenesisState{
				AclAdminList: []types.AclAdmin{
					{
						Address: duplicateAddress,
					},
					{
						Address: duplicateAddress,
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for aclAdmin",
		},
		{
			desc: "invalid aclAdmin address",
			genState: &types.GenesisState{
				AclAdminList: []types.AclAdmin{
					{
						Address: "invalid_address",
					},
					{
						Address: sample.AccAddress(),
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid address for aclAdmin",
		},
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			expErr:   false,
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
				Params: types.DefaultParams(),
				AclAdminList: []types.AclAdmin{
					{
						Address: sample.AccAddress(),
					},
					{
						Address: sample.AccAddress(),
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			expErr: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenesisState_ValidateAclAuthority(t *testing.T) {
	duplicateAddress := sample.AccAddress()
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
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
			},
			expErr:    true,
			expErrMsg: "duplicated index for aclAuthority",
		},
		{
			desc: "invalid aclAuthority address",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: "invalid_address",
					},
					{
						Address: sample.AccAddress(),
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid address for aclAuthority",
		},
		{
			desc: "duplicate access definition module",
			genState: &types.GenesisState{
				AclAuthorityList: []types.AclAuthority{
					{
						Address: sample.AccAddress(),
						Name:    "Alice",
						AccessDefinitions: []*types.AccessDefinition{
							{Module: "module1", IsMaker: true, IsChecker: true},
							{Module: "module1", IsMaker: true, IsChecker: false},
							{Module: "module2", IsMaker: false, IsChecker: true},
						},
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicate module 'module1' found in access definitions",
		},
		{
			desc: "valid aclAuthority",
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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			expErr: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.ValidateAclAuthority()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenesisState_ValidateAclAdmin(t *testing.T) {
	duplicateAddress := sample.AccAddress()
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
	}{
		{
			desc: "duplicated aclAdmin",
			genState: &types.GenesisState{
				AclAdminList: []types.AclAdmin{
					{
						Address: duplicateAddress,
					},
					{
						Address: duplicateAddress,
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for aclAdmin",
		},
		{
			desc: "invalid aclAdmin address",
			genState: &types.GenesisState{
				AclAdminList: []types.AclAdmin{
					{
						Address: "invalid_address",
					},
					{
						Address: sample.AccAddress(),
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid address for aclAdmin",
		},
		{
			desc: "valid aclAdmin",
			genState: &types.GenesisState{
				AclAdminList: []types.AclAdmin{
					{
						Address: sample.AccAddress(),
					},
					{
						Address: sample.AccAddress(),
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			expErr: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.ValidateAclAdmin()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
