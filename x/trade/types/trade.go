package types

// IsValid check if a trade type is valid
func (tt TradeType) IsValid() bool {
	switch tt {
	case TradeTypeBuy,
		TradeTypeSell,
		TradeTypeSplit,
		TradeTypeReinvestment,
		TradeTypeDividends:
		return true
	default:
		return false
	}
}
