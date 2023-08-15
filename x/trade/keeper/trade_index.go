package keeper

import (
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTradeIndex set tradeIndex in the store
func (k Keeper) SetTradeIndex(ctx sdk.Context, tradeIndex types.TradeIndex) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TradeIndexKey))
	b := k.cdc.MustMarshal(&tradeIndex)
	store.Set([]byte{0}, b)
}

// GetTradeIndex returns tradeIndex
func (k Keeper) GetTradeIndex(ctx sdk.Context) (val types.TradeIndex, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TradeIndexKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTradeIndex removes tradeIndex from the store
func (k Keeper) RemoveTradeIndex(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TradeIndexKey))
	store.Delete([]byte{0})
}
