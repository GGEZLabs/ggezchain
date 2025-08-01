package types

const (
	// ModuleName defines the module name
	ModuleName = "trade"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_trade"
)

var ParamsKey = []byte("p_trade")

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TradeIndexKey = "TradeIndex/value/"
)

const (
	DefaultDenom               = "uggz"
	TradeCreatedSuccessfully   = "trade created successfully"
	TradeProcessedSuccessfully = "trade processed successfully"
	TradeIsCanceled            = "trade is canceled"
)

const (
	StatusNil       = TradeStatus_TRADE_STATUS_UNSPECIFIED
	StatusPending   = TradeStatus_TRADE_STATUS_PENDING
	StatusCanceled  = TradeStatus_TRADE_STATUS_CANCELED
	StatusProcessed = TradeStatus_TRADE_STATUS_PROCESSED
	StatusRejected  = TradeStatus_TRADE_STATUS_REJECTED
	StatusFailed    = TradeStatus_TRADE_STATUS_FAILED
)

const (
	ProcessTypeNil     = ProcessType_PROCESS_TYPE_UNSPECIFIED
	ProcessTypeConfirm = ProcessType_PROCESS_TYPE_CONFIRM
	ProcessTypeReject  = ProcessType_PROCESS_TYPE_REJECT
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
