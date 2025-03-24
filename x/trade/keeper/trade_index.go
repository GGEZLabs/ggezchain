package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/x/trade/types"

	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetTradeIndex set tradeIndex in the store
func (k Keeper) SetTradeIndex(ctx context.Context, tradeIndex types.TradeIndex) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TradeIndexKey))
	b := k.cdc.MustMarshal(&tradeIndex)
	store.Set([]byte{0}, b)
}

// GetTradeIndex returns tradeIndex
func (k Keeper) GetTradeIndex(ctx context.Context) (val types.TradeIndex, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TradeIndexKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTradeIndex removes tradeIndex from the store
func (k Keeper) RemoveTradeIndex(ctx context.Context) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TradeIndexKey))
	store.Delete([]byte{0})
}
