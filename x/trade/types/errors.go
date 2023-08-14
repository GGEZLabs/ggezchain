package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/trade module sentinel errors
var (
	ErrTradeCreatedSuccessfully   = sdkerrors.Register(ModuleName, 1000, "trade created successfully")
	ErrTradeProcessedSuccessfully = sdkerrors.Register(ModuleName, 1001, "trade processed successfully")
	ErrSample                     = sdkerrors.Register(ModuleName, 1100, "sample error ")
	ErrInvalidTradeType           = sdkerrors.Register(ModuleName, 1101, "invalid trade type ")
	ErrInvalidTradeData           = sdkerrors.Register(ModuleName, 1102, "invalid trade data ")
	ErrInvalidTradePrice          = sdkerrors.Register(ModuleName, 1103, "invalid trade price ")
	ErrInvalidTradeQuantity       = sdkerrors.Register(ModuleName, 1104, "invalid trade quantity ")
	ErrBurnCoins                  = sdkerrors.Register(ModuleName, 1105, "failed trade burn coins ")
	ErrMintCoins                  = sdkerrors.Register(ModuleName, 1106, "failed trade mint coins ")
	ErrInvalidReceiverAddress     = sdkerrors.Register(ModuleName, 1107, "invalid receiver address ")
	ErrModuleToAccountSendCoins   = sdkerrors.Register(ModuleName, 1108, "failed trade module to account send coin error ")
	ErrAccountToModuleSendCoins   = sdkerrors.Register(ModuleName, 1109, "failed trade account to module send coin error ")
	ErrInvalidCreator             = sdkerrors.Register(ModuleName, 1110, "invalid creator address ")
	ErrInvalidChecker             = sdkerrors.Register(ModuleName, 1111, "invalid checker address ")
	ErrInvalidProcessType         = sdkerrors.Register(ModuleName, 1112, "invalid process type ")
	ErrInvalidTradeIndex          = sdkerrors.Register(ModuleName, 1113, "invalid trade index ")
	ErrTradeStatusCompleted       = sdkerrors.Register(ModuleName, 1114, "trade is already completed ")
	ErrCheckerMustBeDifferent     = sdkerrors.Register(ModuleName, 1115, "checker must be different than maker ")
	ErrTradeStatusCanceled        = sdkerrors.Register(ModuleName, 1116, "trade is already canceled ")
	ErrTradeStatusRejected        = sdkerrors.Register(ModuleName, 1117, "trade is already rejected ")
	ErrInvalidMakerPermission     = sdkerrors.Register(ModuleName, 1118, "invalid maker permission ")
	ErrInvalidCheckerPermission   = sdkerrors.Register(ModuleName, 1119, "invalid checker permission ")
	ErrInvalidDateFormat          = sdkerrors.Register(ModuleName, 1120, "invalid date format ")
	ErrInvalidCreatorAddress      = sdkerrors.Register(ModuleName, 1121, "invalid creator address ")
	ErrInvalidCoinDenom           = sdkerrors.Register(ModuleName, 1122, "invalid coin denom ")
	ErrInvalidPath                = sdkerrors.Register(ModuleName, 1123, "invalid Path ")
	ErrInvalidStatus              = sdkerrors.Register(ModuleName, 1124, "invalid status ")
)
