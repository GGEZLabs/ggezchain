package keeper

import (
	"context"
	"strings"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddAdmin(goCtx context.Context, msg *types.MsgAddAdmin) (*types.MsgAddAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	err := types.ValidateAddAdmin(k.GetAllAclAdmin(ctx), msg.Admins)
	if err != nil {
		return nil, err
	}

	aclAdmins := types.ConvertStringsToAclAdmins(msg.Admins)
	k.SetAclAdmins(ctx, aclAdmins)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddAdmin,
			sdk.NewAttribute(types.AttributeKeyAdmins, strings.Join(msg.Admins, ",")),
		),
	)

	return &types.MsgAddAdminResponse{}, nil
}
