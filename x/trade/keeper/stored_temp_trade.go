package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetStoredTempTrade set a specific storedTempTrade in the store from its index
func (k Keeper) SetStoredTempTrade(ctx context.Context, storedTempTrade types.StoredTempTrade) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTempTradeKeyPrefix))
	b := k.cdc.MustMarshal(&storedTempTrade)
	store.Set(types.StoredTempTradeKey(
		storedTempTrade.TradeIndex,
	), b)
}

// GetStoredTempTrade returns a storedTempTrade from its index
func (k Keeper) GetStoredTempTrade(
	ctx context.Context,
	tradeIndex uint64,

) (val types.StoredTempTrade, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTempTradeKeyPrefix))

	b := store.Get(types.StoredTempTradeKey(
		tradeIndex,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveStoredTempTrade removes a storedTempTrade from the store
func (k Keeper) RemoveStoredTempTrade(
	ctx context.Context,
	tradeIndex uint64,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTempTradeKeyPrefix))
	store.Delete(types.StoredTempTradeKey(
		tradeIndex,
	))
}

// GetAllStoredTempTrade returns all storedTempTrade
func (k Keeper) GetAllStoredTempTrade(ctx context.Context) (list []types.StoredTempTrade) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoredTempTradeKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StoredTempTrade
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
