package keeper

import (
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetStoredTempTrade set a specific storedTempTrade in the store from its index
func (k Keeper) SetStoredTempTrade(ctx sdk.Context, storedTempTrade types.StoredTempTrade) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTempTradeKeyPrefix))
	b := k.cdc.MustMarshal(&storedTempTrade)
	store.Set(types.StoredTempTradeKey(
		storedTempTrade.TradeIndex,
	), b)
}

// GetStoredTempTrade returns a storedTempTrade from its index
func (k Keeper) GetStoredTempTrade(
	ctx sdk.Context,
	tradeIndex uint64,

) (val types.StoredTempTrade, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTempTradeKeyPrefix))

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
	ctx sdk.Context,
	tradeIndex uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTempTradeKeyPrefix))
	store.Delete(types.StoredTempTradeKey(
		tradeIndex,
	))
}

// GetAllStoredTempTrade returns all storedTempTrade
func (k Keeper) GetAllStoredTempTrade(ctx sdk.Context) (list []types.StoredTempTrade) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredTempTradeKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StoredTempTrade
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
