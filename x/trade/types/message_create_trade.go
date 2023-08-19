package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		return ErrInvalidCreatorAddress
	}

	_, err = sdk.AccAddressFromBech32(msg.ReceiverAddress)
	if err != nil {
		return ErrInvalidReceiverAddress
	}

	if msg.Coin != "uggez" {
		return ErrInvalidCoinDenom
	}

	if msg.TradeType != Buy && msg.TradeType != Sell {
		return ErrInvalidTradeType
	}

	if msg.TradeData == "" {
		return ErrInvalidTradeData
	}

	if msg.Quantity == "" {
		return ErrInvalidTradeQuantity
	}

	Quantity, err := strconv.ParseInt(msg.Quantity, 10, 64)
	if err != nil {
		return ErrInvalidTradeQuantity
	}

	if Quantity <= 0 {
		return ErrInvalidTradeQuantity
	}

	if msg.Price == "" {
		return ErrInvalidTradePrice
	}

	CoinPrice, err := strconv.ParseFloat(msg.Price, 32)
	if err != nil {
		return ErrInvalidTradePrice
	}

	if CoinPrice <= 0 {
		return ErrInvalidTradePrice
	}
	return nil
}
