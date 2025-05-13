package types

// DONTCOVER

import (
	sdkErrors "cosmossdk.io/errors"
)

// x/trade module sentinel errors
var (
	ErrInvalidTradeType         = sdkErrors.Register(ModuleName, 1101, "invalid trade type")
	ErrInvalidTradePrice        = sdkErrors.Register(ModuleName, 1102, "invalid trade price")
	ErrInvalidTradeQuantity     = sdkErrors.Register(ModuleName, 1103, "invalid trade quantity")
	ErrInvalidReceiverAddress   = sdkErrors.Register(ModuleName, 1104, "invalid receiver address")
	ErrInvalidProcessType       = sdkErrors.Register(ModuleName, 1105, "invalid process type")
	ErrInvalidTradeIndex        = sdkErrors.Register(ModuleName, 1106, "invalid trade index")
	ErrTradeStatusCompleted     = sdkErrors.Register(ModuleName, 1107, "trade is already completed")
	ErrCheckerMustBeDifferent   = sdkErrors.Register(ModuleName, 1108, "checker must be different than maker")
	ErrTradeStatusCanceled      = sdkErrors.Register(ModuleName, 1109, "trade is already canceled")
	ErrTradeStatusRejected      = sdkErrors.Register(ModuleName, 1110, "trade is already rejected")
	ErrInvalidMakerPermission   = sdkErrors.Register(ModuleName, 1111, "invalid maker permission")
	ErrInvalidCheckerPermission = sdkErrors.Register(ModuleName, 1112, "invalid checker permission")
	ErrInvalidStatus            = sdkErrors.Register(ModuleName, 1113, "invalid status")
	ErrInvalidTradeData         = sdkErrors.Register(ModuleName, 1114, "invalid trade data JSON format")
	ErrInvalidTradeInfo         = sdkErrors.Register(ModuleName, 1115, "invalid trade info")
	ErrInvalidTradeBrokerage    = sdkErrors.Register(ModuleName, 1116, "invalid trade brokerage")
	ErrInvalidSigner            = sdkErrors.Register(ModuleName, 1117, "expected gov account as only signer for proposal message")
	// to check data should be send
	// ErrInvalidCoinMintingPriceJSON = sdkErrors.Register(ModuleName, 1118, "invalid coinMinting price")
	// ErrInvalidExchangeRateJSON     = sdkErrors.Register(ModuleName, 1119, "invalid exchange rate")
	ErrInvalidBankingSystemData = sdkErrors.Register(ModuleName, 1120, "invalid banking system data")
	ErrInvalidMsgType           = sdkErrors.Register(ModuleName, 1121, "invalid message type")
	ErrModuleNotFound           = sdkErrors.Register(ModuleName, 1122, "module not found")
)
