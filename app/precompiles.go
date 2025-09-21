package app

import (
	"fmt"
	"maps"

	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	bankprecompile "github.com/cosmos/evm/precompiles/bank"
	"github.com/cosmos/evm/precompiles/bech32"
	distprecompile "github.com/cosmos/evm/precompiles/distribution"
	// evidenceprecompile "github.com/cosmos/evm/precompiles/evidence"
	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/codec"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	govprecompile "github.com/cosmos/evm/precompiles/gov"
	ics20precompile "github.com/cosmos/evm/precompiles/ics20"
	"github.com/cosmos/evm/precompiles/p256"
	slashingprecompile "github.com/cosmos/evm/precompiles/slashing"
	stakingprecompile "github.com/cosmos/evm/precompiles/staking"
	erc20Keeper "github.com/cosmos/evm/x/erc20/keeper"
	transferkeeper "github.com/cosmos/evm/x/ibc/transfer/keeper"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	channelkeeper "github.com/cosmos/ibc-go/v10/modules/core/04-channel/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const bech32PrecompileBaseGas = 6_000

type PrecompileOptions struct {
	// Codec is the codec used to encode and decode messages.
	AddressCodec address.Codec
	// ValidatorAddressCodec is the codec used to encode and decode validator addresses.
	ValidatorAddressCodec address.Codec
	// ConsensusAddressCodec is the codec used to encode and decode consensus addresses.
	ConsensusAddressCodec address.Codec
}

var CodecOptions = PrecompileOptions{
	AddressCodec:          authcodec.NewBech32Codec("ggez"),
	ValidatorAddressCodec: authcodec.NewBech32Codec("ggezvaloper"),
	ConsensusAddressCodec: authcodec.NewBech32Codec("ggezvalcons"),
}

// TODO: remove this file
// NewAvailableStaticPrecompiles returns all available static precompiled contracts
func NewAvailableStaticPrecompiles(
	cdc codec.Codec,
	stakingKeeper stakingkeeper.Keeper,
	distributionKeeper distributionkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	erc20Keeper erc20Keeper.Keeper,
	authzKeeper authzkeeper.Keeper,
	transferKeeper transferkeeper.Keeper,
	channelKeeper channelkeeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	govKeeper govkeeper.Keeper,
	slashingKeeper slashingkeeper.Keeper,
	// evidenceKeeper evidencekeeper.Keeper,
) map[common.Address]vm.PrecompiledContract {
	precompiles := maps.Clone(vm.PrecompiledContractsBerlin)

	p256Precompile := &p256.Precompile{}

	bech32Precompile, err := bech32.NewPrecompile(bech32PrecompileBaseGas)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate bech32 precompile: %w", err))
	}

	stakingPrecompile, err := stakingprecompile.NewPrecompile(stakingKeeper, CodecOptions.AddressCodec)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate staking precompile: %w", err))
	}

	distributionPrecompile, err := distprecompile.NewPrecompile(distributionKeeper, stakingKeeper, evmKeeper, CodecOptions.AddressCodec)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate distribution precompile: %w", err))
	}

	ibcTransferPrecompile, err := ics20precompile.NewPrecompile(bankKeeper, stakingKeeper, transferKeeper, &channelKeeper, evmKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate ICS20 precompile: %w", err))
	}

	bankPrecompile, err := bankprecompile.NewPrecompile(bankKeeper, erc20Keeper)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate bank precompile: %w", err))
	}

	govPrecompile, err := govprecompile.NewPrecompile(govKeeper, cdc, CodecOptions.AddressCodec)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate gov precompile: %w", err))
	}

	slashingPrecompile, err := slashingprecompile.NewPrecompile(slashingKeeper, CodecOptions.ValidatorAddressCodec, CodecOptions.AddressCodec)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate slashing precompile: %w", err))
	}

	// evidencePrecompile, err := evidenceprecompile.NewPrecompile(evidenceKeeper, authzKeeper)
	// if err != nil {
	//     panic(fmt.Errorf("failed to instantiate evidence precompile: %w", err))
	// }

	// Stateless precompiles
	precompiles[bech32Precompile.Address()] = bech32Precompile
	precompiles[p256Precompile.Address()] = p256Precompile

	// Stateful precompiles
	precompiles[stakingPrecompile.Address()] = stakingPrecompile
	precompiles[distributionPrecompile.Address()] = distributionPrecompile
	precompiles[ibcTransferPrecompile.Address()] = ibcTransferPrecompile
	precompiles[bankPrecompile.Address()] = bankPrecompile
	precompiles[govPrecompile.Address()] = govPrecompile
	precompiles[slashingPrecompile.Address()] = slashingPrecompile
	// precompiles[evidencePrecompile.Address()] = evidencePrecompile

	return precompiles
}
