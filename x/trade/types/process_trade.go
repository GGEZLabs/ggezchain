package types

func (msg *MsgProcessTrade) validateCheckerIsNotMaker(maker string) (err error) {
	if maker == msg.Creator {
		return ErrCheckerMustBeDifferent
	}
	return nil
}

func (msg *MsgProcessTrade) validateStatus(status TradeStatus) (err error) {
	if status != StatusPending {
		return ErrInvalidTradeStatus.Wrapf("cannot process trade with status %s; only trades with status %s can be processed", status.String(), StatusPending.String())
	}
	return nil
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
