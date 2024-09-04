package types

const (
	// ModuleName defines the module name
	ModuleName = "trade"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_trade"
)

var (
	ParamsKey = []byte("p_trade")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TradeIndexKey = "TradeIndex/value/"
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
	DefaultCoinDenom    = "ugz"
	ACLFilePath         = "/.ggezchain/config/chain_acl.json"
)