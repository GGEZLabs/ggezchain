package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Init(ctx context.Context, msg *types.MsgInit) (*types.MsgInitResponse, error) {
	_, err := k.SuperAdmin.Get(ctx)
	if err == nil {
		return nil, types.ErrSuperAdminInitialized
	}
	if !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}

	// Set super admin
	if err := k.SuperAdmin.Set(ctx, types.SuperAdmin{Admin: msg.SuperAdmin}); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInit,
			sdk.NewAttribute(types.AttributeKeySuperAdmin, msg.SuperAdmin),
		),
	)
	return &types.MsgInitResponse{}, nil
}
