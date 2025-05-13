package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgProcessTrade{}

func NewMsgProcessTrade(creator string, processType ProcessType, tradeIndex uint64) *MsgProcessTrade {
	return &MsgProcessTrade{
		Creator:     creator,
		ProcessType: processType,
		TradeIndex:  tradeIndex,
	}
}

func (msg *MsgProcessTrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidCreatorAddress
	}

	if msg.TradeIndex <= 0 {
		return ErrInvalidTradeIndex
	}

	if msg.ProcessType != ProcessTypeConfirm &&
		msg.ProcessType != ProcessTypeReject {
		return ErrInvalidProcessType
	}

	return nil
}
