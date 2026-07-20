package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	authorityStr, err := k.addressCodec.BytesToString(k.GetAuthority())
	if err != nil {
		return nil, err
	}

	if req.Authority != authorityStr {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", authorityStr, req.Authority)
	}

	if err := k.Params.Set(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
