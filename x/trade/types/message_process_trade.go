package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgProcessTrade = "process_trade"

var _ sdk.Msg = &MsgProcessTrade{}

func NewMsgProcessTrade(creator string, processType string, tradeIndex uint64) *MsgProcessTrade {
	return &MsgProcessTrade{
		Creator:     creator,
		ProcessType: processType,
		TradeIndex:  tradeIndex,
	}
}

func (msg *MsgProcessTrade) Route() string {
	return RouterKey
}

func (msg *MsgProcessTrade) Type() string {
	return TypeMsgProcessTrade
}

func (msg *MsgProcessTrade) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgProcessTrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgProcessTrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidCreatorAddress
	}

	if msg.ProcessType != Confirm && msg.ProcessType != Reject {
		return ErrInvalidProcessType
	}

	if msg.TradeIndex <= 0 {
		return ErrInvalidTradeIndex
	}
	return nil
}
