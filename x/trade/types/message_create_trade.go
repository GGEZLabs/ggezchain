package types

import (
	"encoding/json"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.ReceiverAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	if msg.Amount.Denom != DefaultCoinDenom {
		return ErrInvalidDenom
	}

	// todo: check to large quantity
	if msg.Amount.Amount.LTE(math.NewInt(0)) {
		return ErrInvalidTradeQuantity
	}

	if msg.TradeType != TradeTypeBuy &&
		msg.TradeType != TradeTypeSell {
		return ErrInvalidTradeType
	}

	if !json.Valid([]byte(msg.TradeData)) {
		return ErrInvalidTradeData.Wrapf("invalid JSON format for trade-data")
	}

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

	// to check data should be send
	// if msg.CoinMintingPriceJSON == "" {
	// 	return ErrInvalidCoinMintingPriceJSON
	// }

	// if msg.ExchangeRateJSON == "" {
	// 	return ErrInvalidExchangeRateJSON
	// }

	if !json.Valid([]byte(msg.BankingSystemData)) {
		return ErrInvalidBankingSystemData
	}
	return nil
}
