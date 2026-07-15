package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateSuperAdmin(ctx context.Context, msg *types.MsgUpdateSuperAdmin) (*types.MsgUpdateSuperAdminResponse, error) {
	if !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	if err := k.SuperAdmin.Set(ctx, types.SuperAdmin{Admin: msg.NewSuperAdmin}); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateSuperAdmin,
			sdk.NewAttribute(types.AttributeKeySuperAdmin, msg.NewSuperAdmin),
		),
	)
	return &types.MsgUpdateSuperAdminResponse{}, nil
}
