package acl_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/testutil/sample"
	acl "github.com/GGEZLabs/ggezchain/x/acl/module"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AclAuthorityList: []types.AclAuthority{
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
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AclKeeper(t)
	acl.InitGenesis(ctx, k, genesisState)
	got := acl.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.AclAuthorityList, got.AclAuthorityList)
	// this line is used by starport scaffolding # genesis/test/assert
}
