package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ramichain/x/acl/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:          types.DefaultParams(),
		AclAuthorityMap: []types.AclAuthority{{Address: "0"}, {Address: "1"}}, AclAdminMap: []types.AclAdmin{{Address: "0"}, {Address: "1"}}, SuperAdmin: &types.SuperAdmin{Admin: "38"}}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.AclAuthorityMap, got.AclAuthorityMap)
	require.EqualExportedValues(t, genesisState.AclAdminMap, got.AclAdminMap)
	require.EqualExportedValues(t, genesisState.SuperAdmin, got.SuperAdmin)

}
