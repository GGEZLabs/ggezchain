package keeper

import (
	"context"
)

// IsSuperAdmin checks if a given address is a super admin
func (k Keeper) IsSuperAdmin(ctx context.Context, address string) bool {
	superAdmin, found := k.GetSuperAdmin(ctx)

	if found && superAdmin.Admin == address {
		return true
	}

	return false
}

// IsAdmin checks if a given address exists in the list of aclAdmin
func (k Keeper) IsAdmin(ctx context.Context, address string) bool {
	admins := k.GetAllAclAdmin(ctx)

	for _, admin := range admins {
		if admin.Address == address {
			return true
		}
	}

	return false
}
