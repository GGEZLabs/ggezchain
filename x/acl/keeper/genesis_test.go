package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		AclAdmins: []types.AclAdmin{
			{Address: sample.AccAddress()},
			{Address: sample.AccAddress()},
		},
		AclAuthorities: []types.AclAuthority{
			{
				Address: sample.AccAddress(),
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			{
				Address: sample.AccAddress(),
				Name:    "Bob",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
		},
		SuperAdmin: &types.SuperAdmin{Admin: sample.AccAddress()},
	}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.ElementsMatch(t, genesisState.AclAdmins, got.AclAdmins)
	require.ElementsMatch(t, genesisState.AclAuthorities, got.AclAuthorities)
	require.EqualExportedValues(t, genesisState.SuperAdmin, got.SuperAdmin)
}

func TestInitGenesis_InvalidStateRejected(t *testing.T) {
	f := initFixture(t)

	// AclAdmins set without a super admin is invalid per GenesisState.Validate().
	invalidState := types.GenesisState{
		Params: types.DefaultParams(),
		AclAdmins: []types.AclAdmin{
			{Address: sample.AccAddress()},
		},
	}

	err := f.keeper.InitGenesis(f.ctx, invalidState)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot initialize admins or authorities without a super admin")

	// Nothing should have been written to the store.
	_, err = f.keeper.SuperAdmin.Get(f.ctx)
	require.Error(t, err)
}

func TestExportGenesis_NilSuperAdminWhenUnset(t *testing.T) {
	f := initFixture(t)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)

	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.Nil(t, got.SuperAdmin)

	// The exported state must itself be valid, matching what a real
	// export -> re-import (e.g. during a chain upgrade) would require.
	require.NoError(t, got.Validate())
}
