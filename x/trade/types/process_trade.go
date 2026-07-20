package types

const (
	ProcessTypeNil     = ProcessType_PROCESS_TYPE_UNSPECIFIED
	ProcessTypeConfirm = ProcessType_PROCESS_TYPE_CONFIRM
	ProcessTypeReject  = ProcessType_PROCESS_TYPE_REJECT
)

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

// Validate ensures a process-trade request targets a still-pending trade and
// that the checker (msg.Creator) is not the same account that created it.
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
