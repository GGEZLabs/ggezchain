package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteAuthority(goCtx context.Context, msg *types.MsgDeleteAuthority) (*types.MsgDeleteAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAdmin(ctx, msg.Creator) && !k.IsSuperAdmin(ctx, msg.Creator){
		return nil, types.ErrUnauthorized
	}

	_, found := k.GetAclAuthority(ctx, msg.AuthAddress)
	if !found {
		return nil, types.ErrAuthorityAddressDoesNotExist
	}

	k.RemoveAclAuthority(ctx, msg.AuthAddress)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteAuthority,
			sdk.NewAttribute(types.AttributeKeyAuthorityAddress, msg.AuthAddress),
		),
	)

	return &types.MsgDeleteAuthorityResponse{}, nil
}
