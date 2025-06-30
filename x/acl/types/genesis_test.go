package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
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
			desc: "set admins without super admin",
			genState: &types.GenesisState{
				AclAdmins: []types.AclAdmin{
					{
						Address: sample.AccAddress(),
					},
					{
						Address: sample.AccAddress(),
					},
				},
				AclAuthorities: []types.AclAuthority{},
			},
			expErr:    true,
			expErrMsg: "cannot initialize admins or authorities without a super admin",
		},
		{
			desc: "set authorities without super admin",
			genState: &types.GenesisState{
				AclAdmins: []types.AclAdmin{
					{
						Address: sample.AccAddress(),
					},
					{
						Address: sample.AccAddress(),
					},
				},
				AclAuthorities: []types.AclAuthority{},
			},
			expErr:    true,
			expErrMsg: "cannot initialize admins or authorities without a super admin",
		},
		{
			desc: "set authorities without admin",
			genState: &types.GenesisState{
				SuperAdmin: &types.SuperAdmin{
					Admin: sample.AccAddress(),
				},
				AclAuthorities: []types.AclAuthority{
					{
						Address: sample.AccAddress(),
						Name:    "Alice",
					},
				},
			},
			expErr:    true,
			expErrMsg: "cannot initialize authorities without admin",
		},
		{
			desc: "nil super admin",
			genState: &types.GenesisState{
				SuperAdmin: nil,
			},
			expErr: false,
		},
		{
			desc: "invalid super admin address",
			genState: &types.GenesisState{
				SuperAdmin: &types.SuperAdmin{Admin: "invalid_address"},
			},
			expErr:    true,
			expErrMsg: "invalid super admin address",
		},
		{
			desc: "duplicated aclAuthority",
			genState: &types.GenesisState{
				SuperAdmin: &types.SuperAdmin{
					Admin: sample.AccAddress(),
				},
				AclAdmins: []types.AclAdmin{
					{
						Address: sample.AccAddress(),
					},
				},
				AclAuthorities: []types.AclAuthority{
					{
						Address: duplicateAddress,
						Name:    "Alice",
					},
					{
						Address: duplicateAddress,
						Name:    "Bob",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for aclAuthority",
		},
		{
			desc: "invalid aclAuthority address",
			genState: &types.GenesisState{
				SuperAdmin: &types.SuperAdmin{
					Admin: sample.AccAddress(),
				},
				AclAdmins: []types.AclAdmin{
					{
						Address: sample.AccAddress(),
					},
				},
				AclAuthorities: []types.AclAuthority{
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
				SuperAdmin: &types.SuperAdmin{
					Admin: sample.AccAddress(),
				},
				AclAdmins: []types.AclAdmin{
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
				SuperAdmin: &types.SuperAdmin{
					Admin: sample.AccAddress(),
				},
				AclAdmins: []types.AclAdmin{
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
				SuperAdmin: &types.SuperAdmin{
					Admin: sample.AccAddress(),
				},
				AclAuthorities: []types.AclAuthority{
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
				AclAdmins: []types.AclAdmin{
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
				AclAuthorities: []types.AclAuthority{
					{
						Address: duplicateAddress,
						Name:    "Alice",
					},
					{
						Address: duplicateAddress,
						Name:    "Bob",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for aclAuthority",
		},
		{
			desc: "invalid aclAuthority address",
			genState: &types.GenesisState{
				AclAuthorities: []types.AclAuthority{
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
				AclAuthorities: []types.AclAuthority{
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
		// {
		// 	desc: "must have at least one role",
		// 	genState: &types.GenesisState{
		// 		AclAuthorities: []types.AclAuthority{
		// 			{
		// 				Address: sample.AccAddress(),
		// 				Name:    "Alice",
		// 				AccessDefinitions: []*types.AccessDefinition{
		// 					{Module: "module1", IsMaker: false, IsChecker: false},
		// 					{Module: "module2", IsMaker: false, IsChecker: true},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expErr:    true,
		// 	expErrMsg: "access definition for module 'module1' must be either maker or checker",
		// },
		{
			desc: "valid aclAuthority",
			genState: &types.GenesisState{
				AclAuthorities: []types.AclAuthority{
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
				AclAdmins: []types.AclAdmin{
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
				AclAdmins: []types.AclAdmin{
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
				AclAdmins: []types.AclAdmin{
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

func TestGenesisState_ValidateSuperAdmin(t *testing.T) {
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
	}{
		{
			desc: "nil super admin",
			genState: &types.GenesisState{
				SuperAdmin: nil,
			},
			expErr: false,
		},
		{
			desc: "invalid super admin address",
			genState: &types.GenesisState{
				SuperAdmin: &types.SuperAdmin{Admin: "invalid_address"},
			},
			expErr:    true,
			expErrMsg: "invalid super admin address",
		},
		{
			desc: "all good",
			genState: &types.GenesisState{
				SuperAdmin: &types.SuperAdmin{Admin: sample.AccAddress()},
			},
			expErr: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.ValidateSuperAdmin()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
