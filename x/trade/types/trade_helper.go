package types

import (
	_ "github.com/gogo/protobuf/gogoproto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	//sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	errors "cosmossdk.io/errors"
)


func (msg *MsgCreateTrade) ValidateReceiverAddress() (receiverAddress sdk.AccAddress, err error) {
	receiver, errReceiver := sdk.AccAddressFromBech32(msg.ReceiverAddress)
	return receiver, errors.Wrapf(errReceiver, ErrInvalidReceiverAddress.Error())
}

func (msg *MsgCreateTrade) ValidateCreatorAddress() (creatorAddress sdk.AccAddress, err error) {
	creator, errCreator := sdk.AccAddressFromBech32(msg.Creator)
	return creator, errors.Wrapf(errCreator, ErrInvalidCreator.Error())
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

func (msg *MsgProcessTrade) ValidateCheckerAddress(checker string) ( err error) {
	_,errCheckerAddress := sdk.AccAddressFromBech32(checker)
   if errCheckerAddress != nil{
	   return ErrInvalidChecker
   }
   return  nil
}

func (msg *MsgProcessTrade) GetPrepareCoin(storedTrade StoredTrade) (coin sdk.Coin) {
   number, _ := strconv.ParseUint(storedTrade.Quantity, 10, 64)
   return sdk.NewCoin(storedTrade.Coin, sdk.NewInt(int64(number)))
}

func (msg *MsgProcessTrade) checkerAndMakerNotTheSame(maker string ,checker string) (err error) {
   if maker == checker {
	   return ErrCheckerMustBeDifferent
   }
   return nil
}

func (msg *MsgProcessTrade) ValidateStatus(status string) (err error) {
   if status == "Completed" {
	   return  ErrTradeStatusCompleted
   } else if status == "Rejected" {
	   return ErrTradeStatusRejected
   } else if status == "Canceled" {
	   return  ErrTradeStatusCanceled
   }else if status == "Pending" {
	   return nil
   }

   return ErrInvalidStatus
}

/* XXXX func (msg *MsgProcessTrade) CancelExpiredPendingTrades (status string) (err error) {

} */

func (msg *MsgProcessTrade) ValidateProcess(status string, maker string,checker string) (err error) {
	err = msg.ValidateCheckerAddress(checker)
   if err != nil {
	   return err
   }
   err = msg.checkerAndMakerNotTheSame(maker,checker)
   if err != nil {
	   return err
   }
   err = msg.ValidateStatus(status)
   if err != nil {
	   return err
   }

   return err
}

