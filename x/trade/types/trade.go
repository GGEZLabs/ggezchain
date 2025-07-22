package types

// IsTypeValid check if a trade type is valid
func (tt TradeType) IsTypeValid() bool {
	switch tt {
	case TradeTypeBuy,
		TradeTypeSell,
		TradeTypeSplit,
		TradeTypeReverseSplit,
		TradeTypeReinvestment,
		TradeTypeDividends,
		TradeTypeDividendsDeduction:
		return true
	default:
		return false
	}
}

// IsStatusValid check if a trade status is valid
func (ts TradeStatus) IsStatusValid() bool {
	switch ts {
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
