package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateTrade = "create_trade"

var _ sdk.Msg = &MsgCreateTrade{}

func NewMsgCreateTrade(creator string, tradeType string, coin string, price string, quantity string, receiverAddress string, tradeData string) *MsgCreateTrade {
	return &MsgCreateTrade{
		Creator:         creator,
		TradeType:       tradeType,
		Coin:            coin,
		Price:           price,
		Quantity:        quantity,
		ReceiverAddress: receiverAddress,
		TradeData:       tradeData,
	}
}

func (msg *MsgCreateTrade) Route() string {
	return RouterKey
}

func (msg *MsgCreateTrade) Type() string {
	return TypeMsgCreateTrade
}

func (msg *MsgCreateTrade) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
