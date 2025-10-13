package app

import (
	evmconfig "github.com/cosmos/evm/config"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

type EVMOptionsFn func(uint32) error

func NoOpEVMOptions(_ uint32) error {
	return nil
}

// ChainsCoinInfo maps EVM chain IDs to coin configuration
// IMPORTANT: Uses uint64 EVM chain IDs as keys, not Cosmos chain ID strings
var ChainsCoinInfo = map[uint64]evmtypes.EvmCoinInfo{
	EVMChainID: { // Your numeric EVM chain ID (e.g., 9000)
		Denom:         BaseDenom,
		ExtendedDenom: BaseDenom,
		DisplayDenom:  DisplayDenom,
		Decimals:      evmtypes.SixDecimals,
	},
}

// EVMAppOptions sets up global configuration
func EVMAppOptions(chainID uint32) error {
	return evmconfig.EvmAppOptionsWithConfig(EVMChainID, ChainsCoinInfo, cosmosEVMActivators)
}
