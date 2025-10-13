package app

import (
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	evmconfig "github.com/cosmos/evm/config"
	srvflags "github.com/cosmos/evm/server/flags"
	erc20 "github.com/cosmos/evm/x/erc20"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/cosmos/evm/x/feemarket"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	precisebank "github.com/cosmos/evm/x/precisebank"
	precisebankkeeper "github.com/cosmos/evm/x/precisebank/keeper"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	"github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	gethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/spf13/cast"
)

// registerEVMModules register EVM keepers and non dependency inject modules.
func (app *App) registerEVMModules(appOpts servertypes.AppOptions) error {
	displayDenom := "ggez1"
	// chain config
	// TODO:
	coinInfoMap := map[uint64]evmtypes.EvmCoinInfo{
		EVMChainID: {
			Denom:         sdk.DefaultBondDenom,
			ExtendedDenom: sdk.DefaultBondDenom,
			DisplayDenom:  displayDenom,
			Decimals:      evmtypes.EighteenDecimals, // in line with Cosmos SDK default decimals
		},
	}

	// configure evm modules
	if err := evmconfig.EvmAppOptionsWithConfig(
		EVMChainID,
		coinInfoMap,
		getCustomEVMActivators(),
	); err != nil {
		return err
	}

	// set up non depinject support modules store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(evmtypes.StoreKey),
		storetypes.NewKVStoreKey(feemarkettypes.StoreKey),
		storetypes.NewKVStoreKey(erc20types.StoreKey),
		storetypes.NewKVStoreKey(precisebanktypes.StoreKey),
		storetypes.NewKVStoreKey(paramstypes.TStoreKey),
		storetypes.NewTransientStoreKey(evmtypes.TransientKey),
		storetypes.NewTransientStoreKey(feemarkettypes.TransientKey),
	); err != nil {
		return err
	}

	// set up EVM keeper
	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		app.appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.GetKey(feemarkettypes.StoreKey),
		app.UnsafeFindStoreKey(feemarkettypes.TransientKey),
	)

	// PreciseBank wraps BankKeeper to support 18 decimals
	app.PreciseBankKeeper = precisebankkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(precisebanktypes.StoreKey),
		app.BankKeeper,
		app.AuthKeeper,
	)

	// NOTE: it's required to set up the EVM keeper before the ERC-20 keeper, because it is used in its instantiation.
	app.EVMKeeper = evmkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(evmtypes.StoreKey),
		app.UnsafeFindStoreKey(evmtypes.TransientKey),
		app.GetStoreKeysMap(),
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AuthKeeper,
		app.PreciseBankKeeper,
		app.StakingKeeper,
		app.FeeMarketKeeper,
		&app.ConsensusParamsKeeper,
		&app.Erc20Keeper,
		tracer,
	).WithStaticPrecompiles(NewAvailableStaticPrecompiles( // TODO: check precompiles
		*app.StakingKeeper,
		app.DistrKeeper,
		app.BankKeeper,
		app.Erc20Keeper,
		app.TransferKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.EVMKeeper,
		*app.GovKeeper,
		app.SlashingKeeper,
		app.AppCodec(),
	))

	app.Erc20Keeper = erc20keeper.NewKeeper(
		app.GetKey(erc20types.StoreKey),
		app.appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AuthKeeper,
		app.PreciseBankKeeper,
		app.EVMKeeper,
		app.StakingKeeper,
		&app.TransferKeeper,
	)

	// register evm modules
	if err := app.RegisterModules(
		vm.NewAppModule(app.EVMKeeper, app.AuthKeeper, app.AuthKeeper.AddressCodec()),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		erc20.NewAppModule(app.Erc20Keeper, app.AuthKeeper),
		precisebank.NewAppModule(app.PreciseBankKeeper, app.BankKeeper, app.AuthKeeper),
	); err != nil {
		return err
	}

	return nil
}

// RegisterEVM Since the EVM modules don't support dependency injection,
// we need to manually register the modules on the client side.
// This needs to be removed after EVM supports App Wiring.
func RegisterEVM(cdc codec.Codec, interfaceRegistry codectypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		evmtypes.ModuleName:         vm.NewAppModule(nil, authkeeper.AccountKeeper{}, interfaceRegistry.SigningContext().AddressCodec()),
		erc20types.ModuleName:       erc20.NewAppModule(erc20keeper.Keeper{}, authkeeper.AccountKeeper{}),
		feemarkettypes.ModuleName:   feemarket.NewAppModule(feemarketkeeper.Keeper{}),
		precisebanktypes.ModuleName: precisebank.NewAppModule(precisebankkeeper.Keeper{}, bankkeeper.BaseKeeper{}, authkeeper.AccountKeeper{}),
	}

	for _, m := range modules {
		if mr, ok := m.(module.AppModuleBasic); ok {
			mr.RegisterInterfaces(cdc.InterfaceRegistry())
		}
	}

	return modules
}

// ProvideMsgEthereumTxCustomGetSigner provides a custom signer for the MsgEthereumTx message.
func ProvideMsgEthereumTxCustomGetSigner() signing.CustomGetSigner {
	return evmtypes.MsgEthereumTxCustomGetSigner
}

// getCustomEVMActivators defines a map of opcode modifiers associated
// with a key defining the corresponding EIP.
func getCustomEVMActivators() map[int]func(*gethvm.JumpTable) {
	var (
		multiplier        = uint64(10)
		sstoreConstantGas = uint64(500)
	)

	return map[int]func(*gethvm.JumpTable){
		0o000: func(jt *gethvm.JumpTable) {
			// enable0000 contains the logic to modify the CREATE and CREATE2 opcodes
			// constant gas value.
			currentValCreate := jt[gethvm.CREATE].GetConstantGas()
			jt[gethvm.CREATE].SetConstantGas(currentValCreate * multiplier)

			currentValCreate2 := jt[gethvm.CREATE2].GetConstantGas()
			jt[gethvm.CREATE2].SetConstantGas(currentValCreate2 * multiplier)
		},
		0o001: func(jt *gethvm.JumpTable) {
			// enable0001 contains the logic to modify the CALL opcode
			// constant gas value.
			currentVal := jt[gethvm.CALL].GetConstantGas()
			jt[gethvm.CALL].SetConstantGas(currentVal * multiplier)
		},
		0o002: func(jt *gethvm.JumpTable) {
			// enable0002 contains the logic to modify the SSTORE opcode
			// constant gas value.
			jt[gethvm.SSTORE].SetConstantGas(sstoreConstantGas)
		},
	}
}

// GetStoreKeysMap returns a map of store keys.
func (app *App) GetStoreKeysMap() map[string]*storetypes.KVStoreKey {
	storeKeysMap := make(map[string]*storetypes.KVStoreKey)
	for _, storeKey := range app.GetStoreKeys() {
		kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey.Name()).(*storetypes.KVStoreKey)
		if ok {
			storeKeysMap[storeKey.Name()] = kvStoreKey
		}
	}

	return storeKeysMap
}
