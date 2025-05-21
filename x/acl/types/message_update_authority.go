package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateAuthority{}

func NewMsgUpdateAuthority(creator string, authAddress string, newName string, overwriteAccessDefinitions string, addAccessDefinitions string, updateAccessDefinition string, deleteAccessDefinitions []string, clearAllAccessDefinitions bool) *MsgUpdateAuthority {
	return &MsgUpdateAuthority{
		Creator:                    creator,
		AuthAddress:                authAddress,
		NewName:                    newName,
		OverwriteAccessDefinitions: overwriteAccessDefinitions,
		AddAccessDefinitions:       addAccessDefinitions,
		UpdateAccessDefinition:     updateAccessDefinition,
		DeleteAccessDefinitions:    deleteAccessDefinitions,
		ClearAllAccessDefinitions:  clearAllAccessDefinitions,
	}
}

func (msg *MsgUpdateAuthority) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.AuthAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auth-address (%s)", err)
	}

	// Check if none of the flags provided
	hasUpdate := msg.NewName != "" ||
		msg.OverwriteAccessDefinitions != "" ||
		msg.AddAccessDefinitions != "" ||
		msg.UpdateAccessDefinition != "" ||
		len(msg.DeleteAccessDefinitions) > 0 ||
		msg.ClearAllAccessDefinitions

	if !hasUpdate {
		return ErrNoUpdateFlags
	}

	// If OverwriteAccessDefinitions passed ignores other access definition flags
	if msg.OverwriteAccessDefinitions != "" {
		if msg.ClearAllAccessDefinitions || msg.UpdateAccessDefinition != "" || msg.AddAccessDefinitions != "" || len(msg.DeleteAccessDefinitions) > 0 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "overwrite-access-definitions cannot be combined with other access definition flags")
		}
		return validateJSONFormat(msg.OverwriteAccessDefinitions, "overwrite-access-definitions")
	}

	// If ClearAllAccessDefinitions is true ignores other access definition flags
	if msg.ClearAllAccessDefinitions {
		if msg.UpdateAccessDefinition != "" || msg.AddAccessDefinitions != "" || len(msg.DeleteAccessDefinitions) > 0 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "clear-all-access-definitions cannot be combined with other access definition flags")
		}
		return nil
	}

	if msg.UpdateAccessDefinition != "" {
		if err := validateJSONFormat(msg.UpdateAccessDefinition, "update-access-definition"); err != nil {
			return err
		}
	}

	if msg.AddAccessDefinitions != "" {
		if err := validateJSONFormat(msg.AddAccessDefinitions, "add-access-definitions"); err != nil {
			return err
		}
	}

	return nil
}
