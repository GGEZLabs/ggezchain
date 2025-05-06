package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetAclAuthority set a specific aclAuthority in the store from its index
func (k Keeper) SetAclAuthority(ctx context.Context, aclAuthority types.AclAuthority) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAuthorityKeyPrefix))
	b := k.cdc.MustMarshal(&aclAuthority)
	store.Set(types.AclAuthorityKey(
		aclAuthority.Address,
	), b)
}

// GetAclAuthority returns a aclAuthority from its index
func (k Keeper) GetAclAuthority(
	ctx context.Context,
	address string,

) (val types.AclAuthority, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAuthorityKeyPrefix))

	b := store.Get(types.AclAuthorityKey(
		address,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAclAuthority removes a aclAuthority from the store
func (k Keeper) RemoveAclAuthority(
	ctx context.Context,
	address string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAuthorityKeyPrefix))
	store.Delete(types.AclAuthorityKey(
		address,
	))
}

// GetAllAclAuthority returns all aclAuthority
func (k Keeper) GetAllAclAuthority(ctx context.Context) (list []types.AclAuthority) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AclAuthorityKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AclAuthority
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
