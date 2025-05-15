package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/x/acl/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetAclAdmin set a specific aclAdmin in the store from its index
func (k Keeper) SetAclAdmin(ctx context.Context, aclAdmin types.AclAdmin) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAdminKeyPrefix))
	b := k.cdc.MustMarshal(&aclAdmin)
	store.Set(types.AclAdminKey(
		aclAdmin.Address,
	), b)
}

// SetAclAdmins accepts a slice of []types.AclAdmin and adds multiple aclAdmins at once
func (k Keeper) SetAclAdmins(ctx context.Context, aclAdmins []types.AclAdmin) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAdminKeyPrefix))

	for _, aclAdmin := range aclAdmins {
		b := k.cdc.MustMarshal(&aclAdmin)
		store.Set(types.AclAdminKey(aclAdmin.Address), b)
	}
}

// GetAclAdmin returns a aclAdmin from its index
func (k Keeper) GetAclAdmin(
	ctx context.Context,
	address string,
) (val types.AclAdmin, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAdminKeyPrefix))

	b := store.Get(types.AclAdminKey(
		address,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAclAdmin removes a aclAdmin from the store
func (k Keeper) RemoveAclAdmin(
	ctx context.Context,
	address string,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAdminKeyPrefix))
	store.Delete(types.AclAdminKey(
		address,
	))
}

// RemoveAclAdmins accept a []string and remove multiple aclAdmins at once
func (k Keeper) RemoveAclAdmins(ctx context.Context, addresses []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAdminKeyPrefix))

	for _, address := range addresses {
		store.Delete(types.AclAdminKey(address))
	}
}

// GetAllAclAdmin returns all aclAdmin
func (k Keeper) GetAllAclAdmin(ctx context.Context) (list []types.AclAdmin) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAdminKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AclAdmin
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
