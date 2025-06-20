package acl

import (
	"github.com/GGEZLabs/ggezchain/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	err := genState.Validate()
	if err != nil {
		panic(err)
	}

	// Set all the aclAuthority
	for _, elem := range genState.AclAuthorityList {
		k.SetAclAuthority(ctx, elem)
	}
	// Set all the aclAdmin
	for _, elem := range genState.AclAdminList {
		k.SetAclAdmin(ctx, elem)
	}
	// Set if defined
	if genState.SuperAdmin != nil {
		k.SetSuperAdmin(ctx, *genState.SuperAdmin)
	}

	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.AclAuthorityList = k.GetAllAclAuthority(ctx)
	genesis.AclAdminList = k.GetAllAclAdmin(ctx)
	// Get all superAdmin
	superAdmin, found := k.GetSuperAdmin(ctx)
	if found {
		genesis.SuperAdmin = &superAdmin
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
