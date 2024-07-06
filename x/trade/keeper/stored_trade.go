package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetStoredTrade set a specific storedTrade in the store from its index
func (k Keeper) SetStoredTrade(ctx context.Context, storedTrade types.StoredTrade) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTradeKeyPrefix))
	b := k.cdc.MustMarshal(&storedTrade)
	store.Set(types.StoredTradeKey(
		storedTrade.TradeIndex,
	), b)
}

// GetStoredTrade returns a storedTrade from its index
func (k Keeper) GetStoredTrade(
	ctx context.Context,
	tradeIndex uint64,

) (val types.StoredTrade, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTradeKeyPrefix))

	b := store.Get(types.StoredTradeKey(
		tradeIndex,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveStoredTrade removes a storedTrade from the store
func (k Keeper) RemoveStoredTrade(
	ctx context.Context,
	tradeIndex uint64,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTradeKeyPrefix))
	store.Delete(types.StoredTradeKey(
		tradeIndex,
	))
}

// GetAllStoredTrade returns all storedTrade
func (k Keeper) GetAllStoredTrade(ctx context.Context) (list []types.StoredTrade) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTradeKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StoredTrade
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
