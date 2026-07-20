package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	bankKeeper      types.BankKeeper
	aclKeeper       types.AclKeeper
	TradeIndex      collections.Item[types.TradeIndex]
	StoredTrade     collections.Map[uint64, types.StoredTrade]
	StoredTempTrade collections.Map[uint64, types.StoredTempTrade]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	bankKeeper types.BankKeeper,
	aclKeeper types.AclKeeper,
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

		bankKeeper: bankKeeper,
		aclKeeper:  aclKeeper,
		Params:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		TradeIndex: collections.NewItem(sb, types.TradeIndexKey, "tradeIndex", codec.CollValue[types.TradeIndex](cdc)), StoredTrade: collections.NewMap(sb, types.StoredTradeKey, "storedTrade", collections.Uint64Key, codec.CollValue[types.StoredTrade](cdc)), StoredTempTrade: collections.NewMap(sb, types.StoredTempTradeKey, "storedTempTrade", collections.Uint64Key, codec.CollValue[types.StoredTempTrade](cdc)),
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

// AclKeeper exposes the injected AclKeeper for callers outside this package
// (currently only x/trade's simulation operations, which need to seed a sim
// account with trade permissions in x/acl's store before submitting a trade
// message).
func (k Keeper) AclKeeper() types.AclKeeper {
	return k.aclKeeper
}
