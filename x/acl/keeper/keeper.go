package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema       collections.Schema
	Params       collections.Item[types.Params]
	AclAdmin     collections.Map[string, types.AclAdmin]
	AclAuthority collections.Map[string, types.AclAuthority]
	SuperAdmin   collections.Item[types.SuperAdmin]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params:   collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		AclAdmin: collections.NewMap(sb, types.AclAdminKey, "aclAdmin", collections.StringKey, codec.CollValue[types.AclAdmin](cdc)), AclAuthority: collections.NewMap(sb, types.AclAuthorityKey, "aclAuthority", collections.StringKey, codec.CollValue[types.AclAuthority](cdc)), SuperAdmin: collections.NewItem(sb, types.SuperAdminKey, "superAdmin", codec.CollValue[types.SuperAdmin](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// GetAclAuthority returns an AclAuthority by address. Exposed so other modules
// (e.g. x/trade, via its AclKeeper expected-keeper interface) can look up
// authority/permission data without depending on acl's internal storage layout.
func (k Keeper) GetAclAuthority(ctx context.Context, address string) (types.AclAuthority, error) {
	return k.AclAuthority.Get(ctx, address)
}

// SetAclAuthority sets an AclAuthority. Exposed so other modules (e.g. x/trade,
// via its AclKeeper expected-keeper interface) can seed authority/permission
// data — currently only needed by x/trade's simulation operations, which must
// grant a sim account trade permissions before submitting a trade message.
func (k Keeper) SetAclAuthority(ctx context.Context, aclAuthority types.AclAuthority) error {
	return k.AclAuthority.Set(ctx, aclAuthority.Address, aclAuthority)
}
