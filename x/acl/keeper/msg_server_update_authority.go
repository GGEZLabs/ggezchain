package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateAuthority(ctx context.Context, msg *types.MsgUpdateAuthority) (*types.MsgUpdateAuthorityResponse, error) {
	if !k.IsAdmin(ctx, msg.Creator) && !k.IsSuperAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	aclAuthority, err := k.AclAuthority.Get(ctx, msg.AuthAddress)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, types.ErrAuthorityAddressDoesNotExist
		}
		return nil, err
	}

	if msg.NewName != "" {
		aclAuthority = k.UpdateAclAuthorityName(aclAuthority, msg.NewName)
	}

	switch {
	// If OverwriteAccessDefinitions passed ignore another flags
	case msg.OverwriteAccessDefinitions != "":
		aclAuthority, err = k.OverwriteAccessDefinitionList(aclAuthority, msg.OverwriteAccessDefinitions)
		if err != nil {
			return nil, err
		}
	// If ClearAllAccessDefinitions passed ignore another flags
	case msg.ClearAllAccessDefinitions:
		aclAuthority = k.ClearAllAccessDefinitions(aclAuthority)
	default:
		if len(msg.DeleteAccessDefinitions) != 0 {
			if err := types.ValidateDeletedModules(msg.DeleteAccessDefinitions); err != nil {
				return nil, err
			}
		}

		if (msg.UpdateAccessDefinition != "" || msg.AddAccessDefinitions != "") && len(msg.DeleteAccessDefinitions) > 0 {
			if err := types.ValidateConflictBetweenAccessDefinition(msg.UpdateAccessDefinition, msg.AddAccessDefinitions, msg.DeleteAccessDefinitions); err != nil {
				return nil, err
			}
		}

		if msg.UpdateAccessDefinition != "" {
			aclAuthority, err = k.UpdateAccessDefinition(aclAuthority, msg.UpdateAccessDefinition)
			if err != nil {
				return nil, err
			}
		}

		if msg.AddAccessDefinitions != "" {
			aclAuthority, err = k.AddAccessDefinitions(aclAuthority, msg.AddAccessDefinitions)
			if err != nil {
				return nil, err
			}
		}

		if len(msg.DeleteAccessDefinitions) != 0 {
			aclAuthority, err = k.DeleteAccessDefinitions(aclAuthority, msg.DeleteAccessDefinitions)
			if err != nil {
				return nil, err
			}
		}
	}

	// Apply updated aclAuthority
	if err := k.AclAuthority.Set(ctx, aclAuthority.Address, aclAuthority); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateAuthority,
			sdk.NewAttribute(types.AttributeKeyAuthorityAddress, aclAuthority.Address),
			sdk.NewAttribute(types.AttributeKeyName, aclAuthority.Name),
			sdk.NewAttribute(types.AttributeKeyAccessDefinitions, aclAuthority.AccessDefinitionsJson()),
		),
	)

	return &types.MsgUpdateAuthorityResponse{}, nil
}
