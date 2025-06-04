package types_test

import (
	"testing"

	"github.com/GGEZLabs/ramichain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &types.GenesisState{AclAuthorityMap: []types.AclAuthority{{Address: "0"}, {Address: "1"}}, AclAdminMap: []types.AclAdmin{{Address: "0"}, {Address: "1"}}, SuperAdmin: &types.SuperAdmin{Admin: "87"}},
			valid:    true,
		}, {
			desc: "duplicated aclAuthority",
			genState: &types.GenesisState{
				AclAuthorityMap: []types.AclAuthority{
					{
						Address: "0",
					},
					{
						Address: "0",
					},
				},
				AclAdminMap: []types.AclAdmin{{Address: "0"}, {Address: "1"}}, SuperAdmin: &types.SuperAdmin{Admin: "87"}},
			valid: false,
		}, {
			desc: "duplicated aclAdmin",
			genState: &types.GenesisState{
				AclAdminMap: []types.AclAdmin{
					{
						Address: "0",
					},
					{
						Address: "0",
					},
				},
				SuperAdmin: &types.SuperAdmin{Admin: "87"}},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
