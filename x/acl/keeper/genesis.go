package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	for _, elem := range genState.AclAdmins {
		if err := k.AclAdmin.Set(ctx, elem.Address, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.AclAuthorities {
		if err := k.AclAuthority.Set(ctx, elem.Address, elem); err != nil {
			return err
		}
	}
	if genState.SuperAdmin != nil {
		if err := k.SuperAdmin.Set(ctx, *genState.SuperAdmin); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.AclAdmin.Walk(ctx, nil, func(_ string, val types.AclAdmin) (stop bool, err error) {
		genesis.AclAdmins = append(genesis.AclAdmins, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.AclAuthority.Walk(ctx, nil, func(_ string, val types.AclAuthority) (stop bool, err error) {
		genesis.AclAuthorities = append(genesis.AclAuthorities, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	superAdmin, err := k.SuperAdmin.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
	} else {
		genesis.SuperAdmin = &superAdmin
	}

	return genesis, nil
}
