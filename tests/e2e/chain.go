package e2e

import (
	"fmt"
	"os"

	"cosmossdk.io/log"
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	ggezchain "github.com/GGEZLabs/ggezchain/app"
	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"
	tradetypes "github.com/GGEZLabs/ggezchain/x/trade/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsproptypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ratelimittypes "github.com/cosmos/ibc-apps/modules/rate-limiting/v8/types"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var (
	encodingConfig ggezchain.EncodingConfig
	cdc            codec.Codec
	txConfig       client.TxConfig
)

func init() {
	encodingConfig = ggezchain.MakeEncodingConfig()
	banktypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	authvesting.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	stakingtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	evidencetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	govv1types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	govv1beta1types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	paramsproptypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	paramsproptypes.RegisterLegacyAminoCodec(encodingConfig.Amino)
	upgradetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	distribtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ratelimittypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	tradetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	acltypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	cdc = encodingConfig.Codec
	txConfig = encodingConfig.TxConfig
}

type chain struct {
	dataDir    string
	id         string
	validators []*validator
	accounts   []*account //nolint:unused
	// initial accounts in genesis
	genesisAccounts        []*account
	genesisVestingAccounts map[string]sdk.AccAddress
}

func newChain() (*chain, error) {
	tmpDir, err := os.MkdirTemp("", "app-e2e-testnet-")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + tmrand.Str(6),
		dataDir: tmpDir,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int) error {
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = ggezchain.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	tempApplication, _ := ggezchain.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		appOptions,
		baseapp.SetChainID("ggezchain-testnet"),
	)
	defer func() {
		if err := tempApplication.Close(); err != nil {
			panic(err)
		}
	}()

	genesisState := tempApplication.DefaultGenesis()

	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(genesisState); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAndInitValidatorsWithMnemonics(count int, mnemonics []string) error { //nolint:unused // this is called during e2e tests
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = ggezchain.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	tempApplication, _ := ggezchain.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		appOptions,
		baseapp.SetChainID("ggezchain-testnet"),
	)
	defer func() {
		if err := tempApplication.Close(); err != nil {
			panic(err)
		}
	}()

	genesisState := tempApplication.DefaultGenesis()

	for i := 0; i < count; i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(genesisState); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKeyFromMnemonic("val", mnemonics[i]); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createValidator(index int) *validator {
	return &validator{
		chain:   c,
		index:   index,
		moniker: fmt.Sprintf("%s-app-%d", c.id, index),
	}
}
