package app

import (
	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	appante "github.com/GGEZLabs/ggezchain/v2/app/ante"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	evmante "github.com/cosmos/evm/ante"
	cosmosevmante "github.com/cosmos/evm/ante/evm"
	cosmosevmtypes "github.com/cosmos/evm/types"
	"github.com/ethereum/go-ethereum/common"
)

// setAnteHandler sets the ante handler for the application.
func (app *App) setAnteHandler(appCodec codec.Codec, txConfig client.TxConfig, maxGasWanted uint64, wasmConfig wasmtypes.NodeConfig, txCounterStoreKey *storetypes.KVStoreKey) {
	options := appante.HandlerOptions{
		HandlerOptions: evmante.HandlerOptions{
			Cdc:                    appCodec,
			AccountKeeper:          app.AuthKeeper,
			BankKeeper:             app.BankKeeper,
			ExtensionOptionChecker: cosmosevmtypes.HasDynamicFeeExtensionOption,
			EvmKeeper:              app.EVMKeeper,
			FeegrantKeeper:         app.FeeGrantKeeper,
			IBCKeeper:              app.IBCKeeper,
			FeeMarketKeeper:        app.FeeMarketKeeper,
			SignModeHandler:        txConfig.SignModeHandler(),
			SigGasConsumer:         evmante.SigVerificationGasConsumer,
			MaxTxGasWanted:         maxGasWanted,
			TxFeeChecker:           cosmosevmante.NewDynamicFeeChecker(app.FeeMarketKeeper),
			PendingTxListener:      func(hash common.Hash) {},
		},
		IBCKeeper:             app.IBCKeeper,
		NodeConfig:            &wasmConfig,
		WasmKeeper:            &app.WasmKeeper,
		TXCounterStoreService: runtime.NewKVStoreService(txCounterStoreKey),
		CircuitKeeper:         &app.CircuitBreakerKeeper,
	}

	if err := options.Validate(); err != nil {
		panic(err)
	}

	if options.NodeConfig == nil {
		panic("wasm config is required for ante builder")
	}
	if options.TXCounterStoreService == nil {
		panic("wasm store service is required for ante builder")
	}

	app.SetAnteHandler(appante.NewAnteHandler(options))
}
