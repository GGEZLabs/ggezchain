package keeper

import (
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetStoredTrade set a specific storedTrade in the store from its index
func (k Keeper) SetStoredTrade(ctx sdk.Context, storedTrade types.StoredTrade) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTradeKeyPrefix))
	b := k.cdc.MustMarshal(&storedTrade)
	store.Set(types.StoredTradeKey(
		storedTrade.TradeIndex,
	), b)
}

// GetStoredTrade returns a storedTrade from its index
func (k Keeper) GetStoredTrade(
	ctx sdk.Context,
	tradeIndex uint64,

) (val types.StoredTrade, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTradeKeyPrefix))

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
	ctx sdk.Context,
	tradeIndex uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTradeKeyPrefix))
	store.Delete(types.StoredTradeKey(
		tradeIndex,
	))
}

// GetAllStoredTrade returns all storedTrade
func (k Keeper) GetAllStoredTrade(ctx sdk.Context) (list []types.StoredTrade) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTradeKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StoredTrade
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
