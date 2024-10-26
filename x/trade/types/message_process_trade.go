package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgProcessTrade{}

func NewMsgProcessTrade(creator string, processType string, tradeIndex uint64) *MsgProcessTrade {
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

	if msg.ProcessType != Confirm && msg.ProcessType != Reject {
		return ErrInvalidProcessType
	}

	if msg.TradeIndex <= 0 {
		return ErrInvalidTradeIndex
	}
	return nil
}
