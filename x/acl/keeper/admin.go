package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
)

// IsSuperAdmin checks if a given address is a super admin
func (k Keeper) IsSuperAdmin(ctx context.Context, address string) bool {
	superAdmin, err := k.SuperAdmin.Get(ctx)
	if err != nil {
		return false
	}

	return superAdmin.Admin == address
}

// IsAdmin checks if a given address exists in the list of aclAdmin
func (k Keeper) IsAdmin(ctx context.Context, address string) bool {
	found, err := k.AclAdmin.Has(ctx, address)
	if err != nil {
		return false
	}

	return found
}

// GetAllAclAdmin returns all aclAdmin
func (k Keeper) GetAllAclAdmin(ctx context.Context) ([]types.AclAdmin, error) {
	var list []types.AclAdmin

	err := k.AclAdmin.Walk(ctx, nil, func(_ string, val types.AclAdmin) (stop bool, err error) {
		list = append(list, val)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAllAclAuthority returns all aclAuthority
func (k Keeper) GetAllAclAuthority(ctx context.Context) ([]types.AclAuthority, error) {
	var list []types.AclAuthority

	err := k.AclAuthority.Walk(ctx, nil, func(_ string, val types.AclAuthority) (stop bool, err error) {
		list = append(list, val)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
