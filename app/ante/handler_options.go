package ante

import (
	corestoretypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	evmante "github.com/cosmos/evm/ante"
)

// HandlerOptions defines the list of module keepers required to run the EVM
// AnteHandler decorators.
type HandlerOptions struct {
	evmante.HandlerOptions

	NodeConfig             *wasmtypes.NodeConfig
	WasmKeeper             *wasmkeeper.Keeper
	TXCounterStoreService  corestoretypes.KVStoreService
	CircuitKeeper          *circuitkeeper.Keeper
	ExtensionOptionChecker ante.ExtensionOptionChecker
}

// Validate checks if the keepers are defined
func (options HandlerOptions) Validate() error {
	if options.Cdc == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "codec is required for AnteHandler")
	}
	if options.AccountKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SigGasConsumer == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "signature gas consumer is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "sign mode handler is required for AnteHandler")
	}
	if options.CircuitKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "circuit keeper is required for ante builder")
	}

	if options.NodeConfig == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "wasm config is required for ante builder")
	}
	if options.TXCounterStoreService == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "wasm store service is required for ante builder")
	}
	if options.WasmKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "wasm keeper is required for ante builder")
	}

	return nil
}
