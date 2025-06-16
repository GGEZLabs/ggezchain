package keeper

import (
	"context"
	"strings"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddAuthority(goCtx context.Context, msg *types.MsgAddAuthority) (*types.MsgAddAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAdmin(ctx, msg.Creator) && !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	_, found := k.GetAclAuthority(ctx, msg.AuthAddress)
	if found {
		return nil, types.ErrAuthorityAddressExists
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
	k.SetAclAuthority(ctx, aclAuthority)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddAuthority,
			sdk.NewAttribute(types.AttributeKeyAuthorityAddress, aclAuthority.Address),
			sdk.NewAttribute(types.AttributeKeyName, aclAuthority.Name),
			sdk.NewAttribute(types.AttributeKeyAccessDefinitions, aclAuthority.AccessDefinitionsJSON()),
		),
	)
	return &types.MsgAddAuthorityResponse{}, nil
}
