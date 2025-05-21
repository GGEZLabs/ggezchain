package keeper

import (
	"context"
	"strings"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteAclAdmin(goCtx context.Context, msg *types.MsgDeleteAclAdmin) (*types.MsgDeleteAclAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	err := types.ValidateDeleteAclAdmin(k.GetAllAclAdmin(ctx), msg.Admins)
	if err != nil {
		return nil, err
	}

	k.RemoveAclAdmins(ctx, msg.Admins)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteAclAdmin,
			sdk.NewAttribute(types.AttributeKeyAdmins, strings.Join(msg.Admins, ",")),
		),
	)
	return &types.MsgDeleteAclAdminResponse{}, nil
}
