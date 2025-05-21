package keeper

import (
	"context"
	"strings"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) InitAclAdmin(goCtx context.Context, msg *types.MsgInitAclAdmin) (*types.MsgInitAclAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(k.GetAllAclAdmin(ctx)) != 0 {
		return nil, types.ErrAclAdminInitialized
	}

	aclAdmins := types.ConvertStringsToAclAdmins(msg.Admins)
	k.SetAclAdmins(ctx, aclAdmins)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInitAclAdmin,
			sdk.NewAttribute(types.AttributeKeyAdmins, strings.Join(msg.Admins, ",")),
		),
	)
	return &types.MsgInitAclAdminResponse{}, nil
}
