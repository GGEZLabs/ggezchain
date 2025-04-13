package types

import (
	"math"
	"strconv"

	_ "github.com/gogo/protobuf/gogoproto"
	_ "google.golang.org/genproto/googleapis/api/annotations"

	errors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgCreateTrade) ValidateReceiverAndCreatorAddress() (err error) {
	// validate receiver address
	_, err = sdk.AccAddressFromBech32(msg.ReceiverAddress)

	if err != nil {
		return errors.Wrapf(err, ErrInvalidReceiverAddress.Error())
	}
	// validate creator address
	_, err = sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return errors.Wrapf(err, ErrInvalidCreator.Error())
	}
	return nil
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
	return sdk.NewCoin(storedTrade.Coin, sdkmath.NewInt(int64(number))), nil
}

func (msg *MsgProcessTrade) CheckerAndMakerNotTheSame(maker string, checker string) (err error) {
	if maker == checker {
		return ErrCheckerMustBeDifferent
	}
	return nil
}

func (msg *MsgProcessTrade) ValidateStatus(status string) (err error) {
	if status == Completed {
		return ErrTradeStatusCompleted
	} else if status == Rejected {
		return ErrTradeStatusRejected
	} else if status == Canceled {
		return ErrTradeStatusCanceled
	} else if status == Pending {
		return nil
	}

	return ErrInvalidStatus
}

func (msg *MsgProcessTrade) ValidateProcess(status string, maker string, checker string) (err error) {
	err = msg.ValidateCheckerAddress(checker)
	if err != nil {
		return err
	}
	err = msg.CheckerAndMakerNotTheSame(maker, checker)
	if err != nil {
		return err
	}
	err = msg.ValidateStatus(status)
	if err != nil {
		return err
	}

	return err
}
