package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/trade module sentinel errors
var (
	ErrInvalidSigner               = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidTradeType            = errors.Register(ModuleName, 1101, "invalid trade type")
	ErrInvalidTradePrice           = errors.Register(ModuleName, 1102, "invalid trade price")
	ErrInvalidTradeQuantity        = errors.Register(ModuleName, 1103, "invalid trade quantity")
	ErrInvalidReceiverAddress      = errors.Register(ModuleName, 1104, "invalid receiver address")
	ErrInvalidProcessType          = errors.Register(ModuleName, 1105, "invalid process type")
	ErrInvalidTradeIndex           = errors.Register(ModuleName, 1106, "invalid trade index")
	ErrCheckerMustBeDifferent      = errors.Register(ModuleName, 1107, "checker must be different than maker")
	ErrTradeStatusRejected         = errors.Register(ModuleName, 1108, "trade is already rejected")
	ErrInvalidMakerPermission      = errors.Register(ModuleName, 1109, "invalid maker permission")
	ErrInvalidCheckerPermission    = errors.Register(ModuleName, 1110, "invalid checker permission")
	ErrInvalidTradeStatus          = errors.Register(ModuleName, 1111, "invalid status")
	ErrInvalidTradeData            = errors.Register(ModuleName, 1112, "invalid trade data JSON format")
	ErrInvalidTradeInfo            = errors.Register(ModuleName, 1113, "invalid trade info")
	ErrInvalidTradeBrokerage       = errors.Register(ModuleName, 1114, "invalid trade brokerage")
	ErrInvalidCoinMintingPriceJson = errors.Register(ModuleName, 1116, "invalid coin minting price json format")
	ErrInvalidExchangeRateJson     = errors.Register(ModuleName, 1117, "invalid exchange rate json format")
	ErrInvalidBankingSystemData    = errors.Register(ModuleName, 1118, "invalid banking system data json format")
	ErrInvalidMsgType              = errors.Register(ModuleName, 1119, "invalid message type")
	ErrModuleNotFound              = errors.Register(ModuleName, 1120, "module not found")
)
