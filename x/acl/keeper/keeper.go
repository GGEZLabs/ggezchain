package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/GGEZLabs/ramichain/x/acl/types"
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
	AclAuthority collections.Map[string, types.AclAuthority]
	AclAdmin     collections.Map[string, types.AclAdmin]
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

		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		AclAuthority: collections.NewMap(sb, types.AclAuthorityKey, "aclAuthority", collections.StringKey, codec.CollValue[types.AclAuthority](cdc)), AclAdmin: collections.NewMap(sb, types.AclAdminKey, "aclAdmin", collections.StringKey, codec.CollValue[types.AclAdmin](cdc)), SuperAdmin: collections.NewItem(sb, types.SuperAdminKey, "superAdmin", codec.CollValue[types.SuperAdmin](cdc))}

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
