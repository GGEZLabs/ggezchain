package app

import (
    "fmt"

    "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
    evmtypes "github.com/cosmos/evm/x/vm/types"
)

type EVMOptionsFn func(uint32) error

func NoOpEVMOptions(_ uint32) error {
    return nil
}

var sealed = false

// ChainsCoinInfo maps EVM chain IDs to coin configuration
// IMPORTANT: Uses uint64 EVM chain IDs as keys, not Cosmos chain ID strings
var ChainsCoinInfo = map[uint64]evmtypes.EvmCoinInfo{
    EVMChainID: { // Your numeric EVM chain ID (e.g., 9000)
        Denom:        "uggez1",
        ExtendedDenom: "uggez1",
        DisplayDenom: "ggez1",
        Decimals:     evmtypes.EighteenDecimals,
    },
}

// TODO: remove this function
// EVMAppOptions sets up global configuration
func EVMAppOptions(chainID uint32) error {
    if sealed {
        return nil
    }

    // IMPORTANT: Lookup uses numeric EVMChainID, not Cosmos chainID string
    coinInfo, found := ChainsCoinInfo[EVMChainID]
    if !found {
        return fmt.Errorf("unknown EVM chain id: %d", EVMChainID)
    }

    // Set denom info for the chain
    if err := setBaseDenom(coinInfo); err != nil {
        return err
    }

    // baseDenom, err := sdk.GetBaseDenom()
    // if err != nil {
    //     return err
    // }

    ethCfg := evmtypes.DefaultChainConfig(EVMChainID)

    err := evmtypes.NewEVMConfigurator().
        WithChainConfig(ethCfg).
        WithEVMCoinInfo(coinInfo).
        Configure()
    if err != nil {
        return err
    }

    sealed = true
    return nil
}

// setBaseDenom registers display and base denoms
func setBaseDenom(ci evmtypes.EvmCoinInfo) error {
    if err := sdk.RegisterDenom(ci.DisplayDenom, math.LegacyOneDec()); err != nil {
        return err
    }
    return sdk.RegisterDenom(ci.Denom, math.LegacyNewDecWithPrec(1, int64(ci.Decimals)))
}