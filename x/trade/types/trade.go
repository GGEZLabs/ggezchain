package types

const (
	StatusNil       = TradeStatus_TRADE_STATUS_UNSPECIFIED
	StatusPending   = TradeStatus_TRADE_STATUS_PENDING
	StatusCanceled  = TradeStatus_TRADE_STATUS_CANCELED
	StatusProcessed = TradeStatus_TRADE_STATUS_PROCESSED
	StatusRejected  = TradeStatus_TRADE_STATUS_REJECTED
	StatusFailed    = TradeStatus_TRADE_STATUS_FAILED
)

const (
	TradeTypeNil                = TradeType_TRADE_TYPE_UNSPECIFIED
	TradeTypeBuy                = TradeType_TRADE_TYPE_BUY
	TradeTypeSell               = TradeType_TRADE_TYPE_SELL
	TradeTypeSplit              = TradeType_TRADE_TYPE_SPLIT
	TradeTypeReverseSplit       = TradeType_TRADE_TYPE_REVERSE_SPLIT
	TradeTypeReinvestment       = TradeType_TRADE_TYPE_REINVESTMENT
	TradeTypeDividends          = TradeType_TRADE_TYPE_DIVIDENDS
	TradeTypeDividendsDeduction = TradeType_TRADE_TYPE_DIVIDEND_DEDUCTION
)

const (
	TxTypeUnspecified  int32 = 0
	TxTypeCreateTrade  int32 = 1
	TxTypeProcessTrade int32 = 2
)

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
