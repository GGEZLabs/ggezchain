package keeper

import (
	"context"
)

func (k Keeper) IsAdmin(ctx context.Context, address string) bool {
	params := k.GetParams(ctx)
	return params.Admin == address
}
