package keeper

import (
	"context"
	"strings"

	"github.com/GGEZLabs/ggezchain/x/acl/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddAclAdmin(goCtx context.Context, msg *types.MsgAddAclAdmin) (*types.MsgAddAclAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	err := types.ValidateAddAclAdmin(k.GetAllAclAdmin(ctx), msg.Admins)
	if err != nil {
		return nil, err
	}

	aclAdmins := types.ConvertStringsToAclAdmins(msg.Admins)
	k.SetAclAdmins(ctx, aclAdmins)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddAclAdmin,
			sdk.NewAttribute(types.AttributeKeyAdmins, strings.Join(msg.Admins, ",")),
		),
	)

	return &types.MsgAddAclAdminResponse{}, nil
}
