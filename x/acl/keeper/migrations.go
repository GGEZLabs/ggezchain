package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator handles in-place store migrations for the acl module.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator for the acl module.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

const (
	legacySuperAdminKey   = "SuperAdmin/value/"
	legacyAclAdminKey     = "AclAdmin/value/"
	legacyAclAuthorityKey = "AclAuthority/value/"
)

func (m Migrator) MigrateLegacyKeys(ctx sdk.Context) error {
	k := m.keeper
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	// SuperAdmin is a singleton stored under a fixed 0x00 key.
	legacySuperAdmin := prefix.NewStore(storeAdapter, []byte(legacySuperAdminKey))
	if b := legacySuperAdmin.Get([]byte{0}); b != nil {
		var val types.SuperAdmin
		k.cdc.MustUnmarshal(b, &val)
		if err := k.SuperAdmin.Set(ctx, val); err != nil {
			return err
		}
		legacySuperAdmin.Delete([]byte{0})
	}

	if err := migrateLegacyStringMap(storeAdapter, legacyAclAdminKey, func(address string, b []byte) error {
		var val types.AclAdmin
		k.cdc.MustUnmarshal(b, &val)
		return k.AclAdmin.Set(ctx, address, val)
	}); err != nil {
		return err
	}

	if err := migrateLegacyStringMap(storeAdapter, legacyAclAuthorityKey, func(address string, b []byte) error {
		var val types.AclAuthority
		k.cdc.MustUnmarshal(b, &val)
		return k.AclAuthority.Set(ctx, address, val)
	}); err != nil {
		return err
	}

	return nil
}

func migrateLegacyStringMap(storeAdapter storetypes.KVStore, legacyPrefix string, set func(address string, value []byte) error) error {
	store := prefix.NewStore(storeAdapter, []byte(legacyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	var legacyKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if len(key) < 2 || key[len(key)-1] != '/' {
			continue
		}
		address := string(key[:len(key)-1])
		if err := set(address, iterator.Value()); err != nil {
			return err
		}
		legacyKeys = append(legacyKeys, append([]byte{}, key...))
	}

	for _, key := range legacyKeys {
		store.Delete(key)
	}

	return nil
}
