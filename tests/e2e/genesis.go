package e2e

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	tradetypes "github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govmigrv3 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v3"
	govmigrv4 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v4"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govlegacytypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func modifyGenesis(path, moniker, amountStr string, addrAll []sdk.AccAddress, denom string) error {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)
	config.Moniker = moniker

	coins, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}

	var balances []banktypes.Balance
	var genAccounts []*authtypes.BaseAccount
	for _, addr := range addrAll {
		balance := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
		balances = append(balances, balance)
		genAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)
		genAccounts = append(genAccounts, genAccount)
	}

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)
	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return fmt.Errorf("failed to get accounts from any: %w", err)
	}

	for _, addr := range addrAll {
		if accs.Contains(addr) {
			return fmt.Errorf("failed to add account to genesis state; account already exists: %s", addr)
		}
	}

	// Add the new account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	for _, genAcct := range genAccounts {
		accs = append(accs, genAcct)
		accs = authtypes.SanitizeGenesisAccounts(accs)
	}

	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return fmt.Errorf("failed to convert accounts into any's: %w", err)
	}

	authGenState.Accounts = genAccs

	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}
	appState[authtypes.ModuleName] = authGenStateBz

	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances...)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)
	stakingGenState.Params.BondDenom = denom
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal staking genesis state: %s", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// Refactor to separate method
	amnt := math.NewInt(10000)
	quorum, _ := math.LegacyNewDecFromStr("0.000000000000000001")
	threshold, _ := math.LegacyNewDecFromStr("0.000000000000000001")

	maxDepositPeriod := 10 * time.Minute
	votingPeriod := 15 * time.Second
	expeditedVoting := 13 * time.Second

	govStateLegacy := govlegacytypes.NewGenesisState(1,
		govlegacytypes.NewDepositParams(sdk.NewCoins(sdk.NewCoin(denom, amnt)), maxDepositPeriod),
		govlegacytypes.NewVotingParams(votingPeriod),
		govlegacytypes.NewTallyParams(quorum, threshold, govlegacytypes.DefaultVetoThreshold),
	)

	govStateV3, err := govmigrv3.MigrateJSON(govStateLegacy)
	if err != nil {
		return fmt.Errorf("failed to migrate v1beta1 gov genesis state to v3: %w", err)
	}

	govStateV4, err := govmigrv4.MigrateJSON(govStateV3)
	if err != nil {
		return fmt.Errorf("failed to migrate v1beta1 gov genesis state to v4: %w", err)
	}

	govStateV4.Params.ExpeditedVotingPeriod = &expeditedVoting
	govStateV4.Params.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewCoin(denom, amnt)) // same as normal for testing

	govGenStateBz, err := cdc.MarshalJSON(govStateV4)
	if err != nil {
		return fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	tradeGenState := tradetypes.DefaultGenesis()
	tradeGenStateBz, err := cdc.MarshalJSON(tradeGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal trade genesis state: %w", err)
	}
	appState[tradetypes.ModuleName] = tradeGenStateBz

	aclGenState := acltypes.DefaultGenesis()
	aclGenStateBz, err := cdc.MarshalJSON(aclGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal acl genesis state: %w", err)
	}
	appState[acltypes.ModuleName] = aclGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}
	genDoc.AppState = appStateJSON

	return genutil.ExportGenesisFile(genDoc, genFile)
}
