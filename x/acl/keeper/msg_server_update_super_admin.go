package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/GGEZLabs/ramichain/x/acl/types"
)

func (k msgServer) UpdateSuperAdmin(ctx context.Context, msg *types.MsgUpdateSuperAdmin) (*types.MsgUpdateSuperAdminResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgUpdateSuperAdminResponse{}, nil
}
