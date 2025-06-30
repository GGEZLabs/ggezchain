package acl_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/v2/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/v2/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	acl "github.com/GGEZLabs/ggezchain/v2/x/acl/module"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AclAuthorities: []types.AclAuthority{
			{
				Address: sample.AccAddress(),
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
					{Module: "module2", IsMaker: true, IsChecker: true},
				},
			},
			{
				Address: sample.AccAddress(),
				Name:    "Bob",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
					{Module: "module2", IsMaker: true, IsChecker: true},
				},
			},
		},
		AclAdmins: []types.AclAdmin{
			{
				Address: sample.AccAddress(),
			},
			{
				Address: sample.AccAddress(),
			},
		},
		SuperAdmin: &types.SuperAdmin{
			Admin: sample.AccAddress(),
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AclKeeper(t)
	acl.InitGenesis(ctx, k, genesisState)
	got := acl.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.AclAuthorities, got.AclAuthorities)
	require.ElementsMatch(t, genesisState.AclAdmins, got.AclAdmins)
	require.Equal(t, genesisState.SuperAdmin, got.SuperAdmin)
	// this line is used by starport scaffolding # genesis/test/assert
}
