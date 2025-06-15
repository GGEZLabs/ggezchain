package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddAdmin{}

func NewMsgAddAdmin(creator string, admins []string) *MsgAddAdmin {
	return &MsgAddAdmin{
		Creator: creator,
		Admins:  admins,
	}
}

func (msg *MsgAddAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Admins) == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "at least one address should be provided")
	}

	err = validateAddresses(msg.Admins)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address %s", err)
	}

	if hasDuplicateAddresses(msg.Admins) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate address")
	}

	return nil
}
