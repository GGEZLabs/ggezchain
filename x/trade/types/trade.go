package types

// IsTypeValid check if a trade type is valid
func (tt TradeType) IsTypeValid() bool {
	switch tt {
	case TradeTypeBuy,
		TradeTypeSell,
		TradeTypeSplit,
		TradeTypeReinvestment:
		return true
	default:
		return false
	}
}

// IsStatusValid check if a trade status is valid
func (tt TradeStatus) IsStatusValid() bool {
	switch tt {
	case StatusProcessed,
		StatusRejected,
		StatusFailed,
		StatusCanceled,
		StatusPending:
		return true
	default:
		return false
	}
}
