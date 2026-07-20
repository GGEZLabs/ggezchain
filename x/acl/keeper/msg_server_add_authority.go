package keeper

import (
	"context"
	"errors"
	"strings"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddAuthority(ctx context.Context, msg *types.MsgAddAuthority) (*types.MsgAddAuthorityResponse, error) {
	if !k.IsAdmin(ctx, msg.Creator) && !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	_, err := k.AclAuthority.Get(ctx, msg.AuthAddress)
	if err == nil {
		return nil, types.ErrAuthorityAddressExists
	}
	if !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}

	accessDefinitions, err := types.ValidateAccessDefinitionList(msg.AccessDefinitions)
	if err != nil {
		return nil, err
	}

	aclAuthority := types.AclAuthority{
		Address:           msg.AuthAddress,
		Name:              strings.TrimSpace(msg.Name),
		AccessDefinitions: accessDefinitions,
	}

	if err := k.AclAuthority.Set(ctx, aclAuthority.Address, aclAuthority); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddAuthority,
			sdk.NewAttribute(types.AttributeKeyAuthorityAddress, aclAuthority.Address),
			sdk.NewAttribute(types.AttributeKeyName, aclAuthority.Name),
			sdk.NewAttribute(types.AttributeKeyAccessDefinitions, aclAuthority.AccessDefinitionsJson()),
		),
	)
	return &types.MsgAddAuthorityResponse{}, nil
}
