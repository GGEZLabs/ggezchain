package types

import (
	"math"
	"strconv"

	_ "github.com/gogo/protobuf/gogoproto"
	_ "google.golang.org/genproto/googleapis/api/annotations"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errors "cosmossdk.io/errors"
)

func (msg *MsgCreateTrade) ValidateReceiverAddress() (receiverAddress sdk.AccAddress, err error) {
	receiver, err := sdk.AccAddressFromBech32(msg.ReceiverAddress)
	return receiver, errors.Wrapf(err, ErrInvalidReceiverAddress.Error())
}

func (msg *MsgCreateTrade) ValidateCreatorAddress() (creatorAddress sdk.AccAddress, err error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	return creator, errors.Wrapf(err, ErrInvalidCreator.Error())
}

func (msg *MsgCreateTrade) Validate() (err error) {
	_, err = msg.ValidateReceiverAddress()
	if err != nil {
		return err
	}
	_, err = msg.ValidateCreatorAddress()
	if err != nil {
		return err
	}
	return err
}

func (msg *MsgProcessTrade) ValidateCheckerAddress(checker string) (err error) {
	_, err = sdk.AccAddressFromBech32(checker)
	if err != nil {
		return ErrInvalidChecker
	}
	return nil
}

func (msg *MsgProcessTrade) GetPrepareCoin(storedTrade StoredTrade) (sdk.Coin, error) {
	number, err := strconv.ParseUint(storedTrade.Quantity, 10, 64)
	if err != nil {
		return sdk.Coin{}, errors.Wrapf(err, ErrInvalidTradeQuantity.Error())
	}
	if number > math.MaxInt64 {
		return sdk.Coin{}, errors.Wrapf(ErrInvalidTradeQuantity, "quantity too large: %v", storedTrade.Quantity)
	}
	return sdk.NewCoin(storedTrade.Coin, sdk.NewInt(int64(number))), nil
}


func (msg *MsgProcessTrade) checkerAndMakerNotTheSame(maker string, checker string) (err error) {
	if maker == checker {
		return ErrCheckerMustBeDifferent
	}
	return nil
}

func (msg *MsgProcessTrade) ValidateStatus(status string) (err error) {
	if status == "Completed" {
		return ErrTradeStatusCompleted
	} else if status == "Rejected" {
		return ErrTradeStatusRejected
	} else if status == "Canceled" {
		return ErrTradeStatusCanceled
	} else if status == "Pending" {
		return nil
	}

	return ErrInvalidStatus
}

func (msg *MsgProcessTrade) ValidateProcess(status string, maker string, checker string) (err error) {
	err = msg.ValidateCheckerAddress(checker)
	if err != nil {
		return err
	}
	err = msg.checkerAndMakerNotTheSame(maker, checker)
	if err != nil {
		return err
	}
	err = msg.ValidateStatus(status)
	if err != nil {
		return err
	}

	return err
}