package types

// Trade module event types
const (
	EventTypeCreateTrade                = "create_trade"
	EventTypeProcessTrade               = "process_trade"
	EventTypeCancelExpiredPendingTrades = "canceled_trades"

	AttributeKeyTradeIndex  = "trade_index"
	AttributeKeyStatus      = "status"
	AttributeKeyChecker     = "checker"
	AttributeKeyMaker       = "maker"
	AttributeKeyTradeData   = "trade_data"
	AttributeKeyCreateDate  = "create_date"
	AttributeKeyUpdateDate  = "update_date"
	AttributeKeyProcessDate = "process_date"
	AttributeKeyResult      = "result"
)
