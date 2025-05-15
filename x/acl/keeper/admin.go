package keeper

import (
	"context"
)

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
