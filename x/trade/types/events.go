package types

// Trade module event types
const (
	EventTypeCreateTrade                     = "create_trade"
	EventTypeProcessTrade                    = "process_trade"
	EventTypeCanceledTrades                  = "canceled_trades"
	EventTypeCancelExpiredPendingTradesError = "canceled_trades_error"

	AttributeKeyTradeIndex  = "trade_index"
	AttributeKeyStatus      = "status"
	AttributeKeyChecker     = "checker"
	AttributeKeyMaker       = "maker"
	AttributeKeyTradeData   = "trade_data"
	AttributeKeyCreateDate  = "create_date"
	AttributeKeyUpdateDate  = "update_date"
	AttributeKeyProcessDate = "process_date"
	AttributeKeyResult      = "result"
	AttributeKeyError       = "error"
)
