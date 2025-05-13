package types

func (msg *MsgProcessTrade) validateCheckerIsNotMaker(maker string) (err error) {
	if maker == msg.Creator {
		return ErrCheckerMustBeDifferent
	}
	return nil
}

func (msg *MsgProcessTrade) validateStatus(status TradeStatus) (err error) {
	switch status {
	case StatusProcessed:
		return ErrTradeStatusCompleted
	case StatusRejected:
		return ErrTradeStatusRejected
	case StatusCanceled:
		return ErrTradeStatusCanceled
	case StatusPending:
		return nil
	default:
		return ErrInvalidStatus
	}
}

func (msg *MsgProcessTrade) Validate(status TradeStatus, maker string) error {
	err := msg.validateCheckerIsNotMaker(maker)
	if err != nil {
		return err
	}
	err = msg.validateStatus(status)
	if err != nil {
		return err
	}
	return nil
}
