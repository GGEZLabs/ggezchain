package keeper

import (
	"testing"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

// legacyAddressSubKey mirrors the pre-collections key layout: the raw
// address bytes followed by a "/" separator.
func legacyAddressSubKey(address string) []byte {
	return append([]byte(address), '/')
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
	)

	return k, ctx
}

func TestMigrateLegacyKeys(t *testing.T) {
	k, ctx := setupMigrationKeeper(t)

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	legacySuperAdminStore := prefix.NewStore(storeAdapter, []byte(legacySuperAdminKey))
	wantSuperAdmin := types.SuperAdmin{Admin: "legacy-super-admin"}
	legacySuperAdminStore.Set([]byte{0}, k.cdc.MustMarshal(&wantSuperAdmin))

	legacyAclAdminStore := prefix.NewStore(storeAdapter, []byte(legacyAclAdminKey))
	wantAclAdmin := types.AclAdmin{Address: "legacy-admin-addr"}
	legacyAclAdminStore.Set(legacyAddressSubKey(wantAclAdmin.Address), k.cdc.MustMarshal(&wantAclAdmin))

	legacyAclAuthorityStore := prefix.NewStore(storeAdapter, []byte(legacyAclAuthorityKey))
	wantAclAuthority := types.AclAuthority{Address: "legacy-authority-addr", Name: "legacy-authority"}
	legacyAclAuthorityStore.Set(legacyAddressSubKey(wantAclAuthority.Address), k.cdc.MustMarshal(&wantAclAuthority))

	m := NewMigrator(k)
	require.NoError(t, m.MigrateLegacyKeys(ctx))
	// Idempotent: re-running once the legacy keys are gone must be a no-op,
	// since it's registered as the handler for two separate version steps.
	require.NoError(t, m.MigrateLegacyKeys(ctx))

	gotSuperAdmin, err := k.SuperAdmin.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, wantSuperAdmin, gotSuperAdmin)

	gotAclAdmin, err := k.AclAdmin.Get(ctx, wantAclAdmin.Address)
	require.NoError(t, err)
	require.Equal(t, wantAclAdmin, gotAclAdmin)

	gotAclAuthority, err := k.AclAuthority.Get(ctx, wantAclAuthority.Address)
	require.NoError(t, err)
	require.Equal(t, wantAclAuthority, gotAclAuthority)

	require.False(t, legacySuperAdminStore.Has([]byte{0}))
	require.False(t, legacyAclAdminStore.Has(legacyAddressSubKey(wantAclAdmin.Address)))
	require.False(t, legacyAclAuthorityStore.Has(legacyAddressSubKey(wantAclAuthority.Address)))
}
