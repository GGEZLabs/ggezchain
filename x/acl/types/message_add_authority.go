package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddAuthority{}

func NewMsgAddAuthority(creator string, authAddress string, name string, accessDefinitions string) *MsgAddAuthority {
	return &MsgAddAuthority{
		Creator:           creator,
		AuthAddress:       authAddress,
		Name:              name,
		AccessDefinitions: accessDefinitions,
	}
}

func (msg *MsgAddAuthority) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.AuthAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auth-address (%s)", err)
	}

	if msg.Name == "" {
		return ErrEmptyName
	}

	if err = validateJSONFormat(msg.AccessDefinitions, "access-definitions"); err != nil {
		return err
	}

	return nil
}
