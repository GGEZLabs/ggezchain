package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateTrade{}

func NewMsgCreateTrade(creator string, receiverAddress string, tradeData string, bankingSystemData string, coinMintingPriceJson string, exchangeRateJson string) *MsgCreateTrade {
	return &MsgCreateTrade{
		Creator:              creator,
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

	// Validate trade data
	if !json.Valid([]byte(msg.TradeData)) {
		return ErrInvalidTradeData
	}

	// Validate banking system data
	if !json.Valid([]byte(msg.BankingSystemData)) {
		return ErrInvalidBankingSystemData
	}

	// Validate CoinMintingPriceJson
	if !json.Valid([]byte(msg.CoinMintingPriceJson)) {
		return ErrInvalidCoinMintingPriceJson
	}

	// Validate ExchangeRateJson
	if !json.Valid([]byte(msg.ExchangeRateJson)) {
		return ErrInvalidExchangeRateJson
	}

	// Validate CreateDate if it does not empty
	if msg.CreateDate != "" {
		_, err = time.Parse(time.RFC3339, msg.CreateDate)
		if err != nil {
			return sdkerrors.ErrInvalidRequest.Wrapf("invalid create_date format: %s, date format should be like: %s", msg.CreateDate, time.RFC3339)
		}
	}

	return nil
}
