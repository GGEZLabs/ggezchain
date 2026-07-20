package keeper

import (
	"context"
	"strings"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteAdmin(ctx context.Context, msg *types.MsgDeleteAdmin) (*types.MsgDeleteAdminResponse, error) {
	if !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	currentAdmins, err := k.GetAllAclAdmin(ctx)
	if err != nil {
		return nil, err
	}

	if err := types.ValidateDeleteAdmin(currentAdmins, msg.Admins); err != nil {
		return nil, err
	}

	for _, address := range msg.Admins {
		if err := k.AclAdmin.Remove(ctx, address); err != nil {
			return nil, err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeleteAdmin,
			sdk.NewAttribute(types.AttributeKeyAdmins, strings.Join(msg.Admins, ",")),
		),
	)

	return &types.MsgDeleteAdminResponse{}, nil
}
