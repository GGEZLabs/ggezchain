package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteAuthority(ctx context.Context, msg *types.MsgDeleteAuthority) (*types.MsgDeleteAuthorityResponse, error) {
	if !k.IsAdmin(ctx, msg.Creator) && !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	_, err := k.AclAuthority.Get(ctx, msg.AuthAddress)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, types.ErrAuthorityAddressDoesNotExist
		}
		return nil, err
	}

	if err := k.AclAuthority.Remove(ctx, msg.AuthAddress); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteAuthority,
			sdk.NewAttribute(types.AttributeKeyAuthorityAddress, msg.AuthAddress),
		),
	)

	return &types.MsgDeleteAuthorityResponse{}, nil
}
