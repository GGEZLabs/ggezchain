package types

const (
	// ModuleName defines the module name
	ModuleName = "trade"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_trade"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TradeIndexKey = "TradeIndex/value/"
)
const (
	CancelExpiredPendingTradesEventType = "status is changed"
)
const (
	Pending             = "Pending"
	Failed              = "Failed"
	Rejected            = "Rejected"
	Canceled            = "Canceled"
	Completed           = "Completed"
	CoinsStuckOnModule  = "Coins Stuck On Module"
	CoinsStuckOnAccount = "Coins Stuck On Account"
	Confirm             = "Confirm"
	Reject              = "Reject"
	Buy                 = "buy"
	Sell                = "sell"
	CreateTrade         = "CreateTrade"
	ProcessTrade        = "ProcessTrade"
	DefaultCoinDenom    = "uggez"
)
