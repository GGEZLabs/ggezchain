package types

import (
	"encoding/json"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateTrade{}

func NewMsgCreateTrade(creator string, tradeType TradeType, amount *sdk.Coin, price string, receiverAddress string, tradeData string, bankingSystemData string, coinMintingPriceJson string, exchangeRateJson string) *MsgCreateTrade {
	return &MsgCreateTrade{
		Creator:              creator,
		TradeType:            tradeType,
		Amount:               amount,
		Price:                price,
		ReceiverAddress:      receiverAddress,
		TradeData:            tradeData,
		BankingSystemData:    bankingSystemData,
		CoinMintingPriceJson: coinMintingPriceJson,
		ExchangeRateJson:     exchangeRateJson,
	}
}

func (msg *MsgCreateTrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid creator address (%s)", err)
	}

	// Validate trade type
	if msg.TradeType != TradeTypeBuy &&
		msg.TradeType != TradeTypeSell {
		return ErrInvalidTradeType
	}

	// Validate amount
	if !msg.Amount.IsValid() {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid amount: %s", msg.Amount.String())
	}

	if msg.Amount.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrapf("zero amount not allowed: %s", msg.Amount.String())
	}

	if msg.Amount.Denom != DefaultDenom {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid denom expected: %s, got: %s ", DefaultDenom, msg.Amount.Denom)
	}

	// Validate price
	if strings.TrimSpace(msg.Price) == "" {
		return ErrInvalidTradePrice
	}

	coinPrice, err := strconv.ParseFloat(msg.Price, 64)
	if err != nil {
		return ErrInvalidTradePrice.Wrapf(err.Error())
	}

	if coinPrice <= 0 {
		return ErrInvalidTradePrice
	}

	// Validate receiver address
	_, err = sdk.AccAddressFromBech32(msg.ReceiverAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address (%s)", err)
	}

	// Validate trade data
	if !json.Valid([]byte(msg.TradeData)) {
		return ErrInvalidTradeData
	}

	// Validate banking system data
	if !json.Valid([]byte(msg.BankingSystemData)) {
		return ErrInvalidBankingSystemData
	}

	// to check data should be send
	// if msg.CoinMintingPriceJSON == "" {
	// 	return ErrInvalidCoinMintingPriceJSON
	// }

	// if msg.ExchangeRateJSON == "" {
	// 	return ErrInvalidExchangeRateJSON
	// }

	return nil
}
