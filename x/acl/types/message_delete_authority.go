package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgDeleteAuthority{}

func NewMsgDeleteAuthority(creator string, authAddress string) *MsgDeleteAuthority {
	return &MsgDeleteAuthority{
		Creator:     creator,
		AuthAddress: authAddress,
	}
}

func (msg *MsgDeleteAuthority) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.AuthAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auth-address (%s)", err)
	}

	return nil
}
