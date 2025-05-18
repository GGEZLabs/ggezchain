package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/trade module sentinel errors
var (
	ErrInvalidTradeType         = sdkerrors.Register(ModuleName, 1101, "invalid trade type")
	ErrInvalidTradePrice        = sdkerrors.Register(ModuleName, 1102, "invalid trade price")
	ErrInvalidTradeQuantity     = sdkerrors.Register(ModuleName, 1103, "invalid trade quantity")
	ErrInvalidReceiverAddress   = sdkerrors.Register(ModuleName, 1104, "invalid receiver address")
	ErrInvalidProcessType       = sdkerrors.Register(ModuleName, 1105, "invalid process type")
	ErrInvalidTradeIndex        = sdkerrors.Register(ModuleName, 1106, "invalid trade index")
	ErrCheckerMustBeDifferent   = sdkerrors.Register(ModuleName, 1107, "checker must be different than maker")
	ErrTradeStatusCompleted     = sdkerrors.Register(ModuleName, 1108, "trade is already completed")
	ErrTradeStatusCanceled      = sdkerrors.Register(ModuleName, 1109, "trade is already canceled")
	ErrTradeStatusRejected      = sdkerrors.Register(ModuleName, 1110, "trade is already rejected")
	ErrInvalidMakerPermission   = sdkerrors.Register(ModuleName, 1111, "invalid maker permission")
	ErrInvalidCheckerPermission = sdkerrors.Register(ModuleName, 1112, "invalid checker permission")
	ErrInvalidStatus            = sdkerrors.Register(ModuleName, 1113, "invalid status")
	ErrInvalidTradeData         = sdkerrors.Register(ModuleName, 1114, "invalid trade data JSON format")
	ErrInvalidTradeInfo         = sdkerrors.Register(ModuleName, 1115, "invalid trade info")
	ErrInvalidTradeBrokerage    = sdkerrors.Register(ModuleName, 1116, "invalid trade brokerage")
	ErrInvalidSigner            = sdkerrors.Register(ModuleName, 1117, "expected gov account as only signer for proposal message")
	// todo: to check data should be send
	// ErrInvalidCoinMintingPriceJSON = sdkErrors.Register(ModuleName, 1118, "invalid coinMinting price")
	// ErrInvalidExchangeRateJSON     = sdkErrors.Register(ModuleName, 1119, "invalid exchange rate")
	ErrInvalidBankingSystemData = sdkerrors.Register(ModuleName, 1120, "invalid banking system data")
	ErrInvalidMsgType           = sdkerrors.Register(ModuleName, 1121, "invalid message type")
	ErrModuleNotFound           = sdkerrors.Register(ModuleName, 1122, "module not found")
)
