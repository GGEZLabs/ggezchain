package keeper

import (
	"context"
	"strings"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddAdmin(ctx context.Context, msg *types.MsgAddAdmin) (*types.MsgAddAdminResponse, error) {
	if !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	currentAdmins, err := k.GetAllAclAdmin(ctx)
	if err != nil {
		return nil, err
	}

	if err := types.ValidateAddAdmin(currentAdmins, msg.Admins); err != nil {
		return nil, err
	}

	aclAdmins := types.ConvertStringsToAclAdmins(msg.Admins)
	for _, aclAdmin := range aclAdmins {
		if err := k.AclAdmin.Set(ctx, aclAdmin.Address, aclAdmin); err != nil {
			return nil, err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddAdmin,
			sdk.NewAttribute(types.AttributeKeyAdmins, strings.Join(msg.Admins, ",")),
		),
	)

	return &types.MsgAddAdminResponse{}, nil
}
