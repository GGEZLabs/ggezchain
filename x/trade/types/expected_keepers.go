package types

import (
	"context"

	"cosmossdk.io/core/address"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AclKeeper interface {
	GetAclAuthority(ctx context.Context, address string) (acltypes.AclAuthority, error)
	// SetAclAuthority is only used by simulation, to seed a sim account with
	// trade permissions before submitting a trade message.
	SetAclAuthority(ctx context.Context, aclAuthority acltypes.AclAuthority) error
}

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
