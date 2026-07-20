package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator handles in-place store migrations for the trade module.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator for the trade module.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

const (
	legacyTradeIndexKey      = "TradeIndex/value/"
	legacyStoredTradeKey     = "StoredTrade/value/"
	legacyStoredTempTradeKey = "StoredTempTrade/value/"
)

func (m Migrator) MigrateLegacyKeys(ctx sdk.Context) error {
	k := m.keeper
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	// TradeIndex is a singleton stored under a fixed 0x00 key.
	legacyTradeIndex := prefix.NewStore(storeAdapter, []byte(legacyTradeIndexKey))
	if b := legacyTradeIndex.Get([]byte{0}); b != nil {
		var val types.TradeIndex
		k.cdc.MustUnmarshal(b, &val)
		if err := k.TradeIndex.Set(ctx, val); err != nil {
			return err
		}
		legacyTradeIndex.Delete([]byte{0})
	}

	if err := migrateLegacyUint64Map(storeAdapter, legacyStoredTradeKey, func(idx uint64, b []byte) error {
		var val types.StoredTrade
		k.cdc.MustUnmarshal(b, &val)
		return k.StoredTrade.Set(ctx, idx, val)
	}); err != nil {
		return err
	}

	if err := migrateLegacyUint64Map(storeAdapter, legacyStoredTempTradeKey, func(idx uint64, b []byte) error {
		var val types.StoredTempTrade
		k.cdc.MustUnmarshal(b, &val)
		return k.StoredTempTrade.Set(ctx, idx, val)
	}); err != nil {
		return err
	}

	return nil
}

func migrateLegacyUint64Map(storeAdapter storetypes.KVStore, legacyPrefix string, set func(idx uint64, value []byte) error) error {
	store := prefix.NewStore(storeAdapter, []byte(legacyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	var legacyKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if len(key) < 8 {
			continue
		}
		idx := binary.BigEndian.Uint64(key[:8])
		if err := set(idx, iterator.Value()); err != nil {
			return err
		}
		legacyKeys = append(legacyKeys, append([]byte{}, key...))
	}

	for _, key := range legacyKeys {
		store.Delete(key)
	}

	return nil
}
