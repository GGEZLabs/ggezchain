package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateAuthority(goCtx context.Context, msg *types.MsgUpdateAuthority) (*types.MsgUpdateAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAdmin(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	aclAuthority, found := k.GetAclAuthority(ctx, msg.AuthAddress)
	if !found {
		return nil, types.ErrAuthorityAddressNotExist
	}

	if msg.NewName != "" {
		aclAuthority = k.UpdateAclAuthorityName(aclAuthority, msg.NewName)
	}

	var err error
	// if OverwriteAccessDefinitions passed ignore another flags
	if msg.OverwriteAccessDefinitions != "" {
		aclAuthority, err = k.OverwriteAccessDefinitionList(aclAuthority, msg.OverwriteAccessDefinitions)
		if err != nil {
			return nil, err
		}
	} else if msg.ClearAllAccessDefinitions {
		// if ClearAllAccessDefinitions passed ignore another flags
		aclAuthority = k.ClearAllAccessDefinitions(aclAuthority)
	} else {

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
			aclAuthority, err = k.UpdateAccessDefinitions(aclAuthority, msg.UpdateAccessDefinition)
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
	// apply updated aclAuthority
	k.SetAclAuthority(ctx, aclAuthority)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateAuthority,
			sdk.NewAttribute(types.AttributeKeyAuthorityAddress, aclAuthority.Address),
			sdk.NewAttribute(types.AttributeKeyName, aclAuthority.Name),
			sdk.NewAttribute(types.AttributeKeyAccessDefinitions, aclAuthority.AccessDefinitionsJSON()),
		),
	)

	return &types.MsgUpdateAuthorityResponse{}, nil
}
