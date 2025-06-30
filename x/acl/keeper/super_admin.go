package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetSuperAdmin set superAdmin in the store
func (k Keeper) SetSuperAdmin(ctx context.Context, superAdmin types.SuperAdmin) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SuperAdminKey))
	b := k.cdc.MustMarshal(&superAdmin)
	store.Set([]byte{0}, b)
}

// GetSuperAdmin returns superAdmin
func (k Keeper) GetSuperAdmin(ctx context.Context) (val types.SuperAdmin, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SuperAdminKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSuperAdmin removes superAdmin from the store
func (k Keeper) RemoveSuperAdmin(ctx context.Context) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SuperAdminKey))
	store.Delete([]byte{0})
}
