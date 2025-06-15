package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateSuperAdmin{}

func NewMsgUpdateSuperAdmin(creator string, newSuperAdmin string) *MsgUpdateSuperAdmin {
	return &MsgUpdateSuperAdmin{
		Creator:       creator,
		NewSuperAdmin: newSuperAdmin,
	}
}

func (msg *MsgUpdateSuperAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.NewSuperAdmin)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new-super-admin address (%s)", err)
	}

	return nil
}
