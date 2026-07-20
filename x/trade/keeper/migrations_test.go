package keeper

import (
	"encoding/binary"
	"testing"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

// legacyTradeSubKey mirrors the pre-collections key layout: an 8-byte
// big-endian index followed by a "/" separator.
func legacyTradeSubKey(idx uint64) []byte {
	b := binary.BigEndian.AppendUint64(make([]byte, 0, 9), idx)
	return append(b, '/')
}

func setupMigrationKeeper(t *testing.T) (Keeper, sdk.Context) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig()
	addrCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	k := NewKeeper(
		storeService,
		encCfg.Codec,
		addrCodec,
		authtypes.NewModuleAddress(types.GovModuleName),
		nil,
		nil,
	)

	return k, ctx
}

func TestMigrateLegacyKeys(t *testing.T) {
	k, ctx := setupMigrationKeeper(t)

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	legacyTradeIndexStore := prefix.NewStore(storeAdapter, []byte(legacyTradeIndexKey))
	wantTradeIndex := types.TradeIndex{NextId: 42}
	legacyTradeIndexStore.Set([]byte{0}, k.cdc.MustMarshal(&wantTradeIndex))

	legacyStoredTradeStore := prefix.NewStore(storeAdapter, []byte(legacyStoredTradeKey))
	wantStoredTrade := types.StoredTrade{TradeIndex: 7, Maker: "legacy-maker"}
	legacyStoredTradeStore.Set(legacyTradeSubKey(7), k.cdc.MustMarshal(&wantStoredTrade))

	legacyStoredTempTradeStore := prefix.NewStore(storeAdapter, []byte(legacyStoredTempTradeKey))
	wantStoredTempTrade := types.StoredTempTrade{TradeIndex: 9, TxDate: "legacy-temp-tx-date"}
	legacyStoredTempTradeStore.Set(legacyTradeSubKey(9), k.cdc.MustMarshal(&wantStoredTempTrade))

	m := NewMigrator(k)
	require.NoError(t, m.MigrateLegacyKeys(ctx))
	// Idempotent: re-running once the legacy keys are gone must be a no-op,
	// since it's registered as the handler for two separate version steps.
	require.NoError(t, m.MigrateLegacyKeys(ctx))

	gotTradeIndex, err := k.TradeIndex.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, wantTradeIndex, gotTradeIndex)

	gotStoredTrade, err := k.StoredTrade.Get(ctx, 7)
	require.NoError(t, err)
	require.Equal(t, wantStoredTrade, gotStoredTrade)

	gotStoredTempTrade, err := k.StoredTempTrade.Get(ctx, 9)
	require.NoError(t, err)
	require.Equal(t, wantStoredTempTrade, gotStoredTempTrade)

	require.False(t, legacyTradeIndexStore.Has([]byte{0}))
	require.False(t, legacyStoredTradeStore.Has(legacyTradeSubKey(7)))
	require.False(t, legacyStoredTempTradeStore.Has(legacyTradeSubKey(9)))
}
