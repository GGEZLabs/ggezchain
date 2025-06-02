package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cosmossdk.io/math"
	evidencetypes "cosmossdk.io/x/evidence/types"

	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"
	tradetypes "github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ratelimittypes "github.com/cosmos/ibc-apps/modules/rate-limiting/v8/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func queryTx(endpoint, txHash string) error {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", endpoint, txHash))
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tx query returned non-200 status: %d", resp.StatusCode)
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	txResp := result["tx_response"].(map[string]interface{})
	if v := txResp["code"]; v.(float64) != 0 {
		return fmt.Errorf("tx %s failed with status code %v", txHash, v)
	}

	return nil
}

// if coin is zero, return empty coin.
func getSpecificBalance(endpoint, addr, denom string) (amt sdk.Coin, err error) {
	balances, err := queryAllBalances(endpoint, addr)
	if err != nil {
		return amt, err
	}
	amt = sdk.NewCoin(denom, math.ZeroInt())
	for _, c := range balances {
		if strings.Contains(c.Denom, denom) {
			amt = c
			break
		}
	}
	return amt, nil
}

func queryAllBalances(endpoint, addr string) (sdk.Coins, error) {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", endpoint, addr))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var balancesResp banktypes.QueryAllBalancesResponse
	if err := cdc.UnmarshalJSON(body, &balancesResp); err != nil {
		return nil, err
	}

	return balancesResp.Balances, nil
}

func querySupplyOf(endpoint, denom string) (sdk.Coin, error) {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/bank/v1beta1/supply/by_denom?denom=%s", endpoint, denom))
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var supplyOfResp banktypes.QuerySupplyOfResponse
	if err := cdc.UnmarshalJSON(body, &supplyOfResp); err != nil {
		return sdk.Coin{}, err
	}

	return supplyOfResp.Amount, nil
}

// func queryStakingParams(endpoint string) (stakingtypes.QueryParamsResponse, error) {
// 	body, err := httpGet(fmt.Sprintf("%s/cosmos/staking/v1beta1/params", endpoint))
// 	if err != nil {
// 		return stakingtypes.QueryParamsResponse{}, fmt.Errorf("failed to execute HTTP request: %w", err)
// 	}

// 	var params stakingtypes.QueryParamsResponse
// 	if err := cdc.UnmarshalJSON(body, &params); err != nil {
// 		return stakingtypes.QueryParamsResponse{}, err
// 	}

// 	return params, nil
// }

func queryDelegation(endpoint string, validatorAddr string, delegatorAddr string) (stakingtypes.QueryDelegationResponse, error) {
	var res stakingtypes.QueryDelegationResponse

	body, err := httpGet(fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations/%s", endpoint, validatorAddr, delegatorAddr))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryUnbondingDelegation(endpoint string, validatorAddr string, delegatorAddr string) (stakingtypes.QueryUnbondingDelegationResponse, error) {
	var res stakingtypes.QueryUnbondingDelegationResponse
	body, err := httpGet(fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations/%s/unbonding_delegation", endpoint, validatorAddr, delegatorAddr))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryDelegatorWithdrawalAddress(endpoint string, delegatorAddr string) (disttypes.QueryDelegatorWithdrawAddressResponse, error) {
	var res disttypes.QueryDelegatorWithdrawAddressResponse

	body, err := httpGet(fmt.Sprintf("%s/cosmos/distribution/v1beta1/delegators/%s/withdraw_address", endpoint, delegatorAddr))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryGovProposal(endpoint string, proposalID int) (govtypesv1beta1.QueryProposalResponse, error) {
	var govProposalResp govtypesv1beta1.QueryProposalResponse

	path := fmt.Sprintf("%s/cosmos/gov/v1beta1/proposals/%d", endpoint, proposalID)

	body, err := httpGet(path)
	if err != nil {
		return govProposalResp, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	if err := cdc.UnmarshalJSON(body, &govProposalResp); err != nil {
		return govProposalResp, err
	}

	return govProposalResp, nil
}

func queryGovProposalV1(endpoint string, proposalID int) (govtypesv1.QueryProposalResponse, error) {
	var govProposalResp govtypesv1.QueryProposalResponse

	path := fmt.Sprintf("%s/cosmos/gov/v1/proposals/%d", endpoint, proposalID)

	body, err := httpGet(path)
	if err != nil {
		return govProposalResp, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	if err := cdc.UnmarshalJSON(body, &govProposalResp); err != nil {
		return govProposalResp, err
	}

	return govProposalResp, nil
}

func queryAccount(endpoint, address string) (acc sdk.AccountI, err error) {
	var res authtypes.QueryAccountResponse
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", endpoint, address))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := cdc.UnmarshalJSON(bz, &res); err != nil {
		return nil, err
	}
	return acc, cdc.UnpackAny(res.Account, &acc)
}

func queryDelayedVestingAccount(endpoint, address string) (authvesting.DelayedVestingAccount, error) {
	baseAcc, err := queryAccount(endpoint, address)
	if err != nil {
		return authvesting.DelayedVestingAccount{}, err
	}
	acc, ok := baseAcc.(*authvesting.DelayedVestingAccount)
	if !ok {
		return authvesting.DelayedVestingAccount{},
			fmt.Errorf("cannot cast %v to DelayedVestingAccount", baseAcc)
	}
	return *acc, nil
}

func queryContinuousVestingAccount(endpoint, address string) (authvesting.ContinuousVestingAccount, error) {
	baseAcc, err := queryAccount(endpoint, address)
	if err != nil {
		return authvesting.ContinuousVestingAccount{}, err
	}
	acc, ok := baseAcc.(*authvesting.ContinuousVestingAccount)
	if !ok {
		return authvesting.ContinuousVestingAccount{},
			fmt.Errorf("cannot cast %v to ContinuousVestingAccount", baseAcc)
	}
	return *acc, nil
}

func queryPeriodicVestingAccount(endpoint, address string) (authvesting.PeriodicVestingAccount, error) {
	baseAcc, err := queryAccount(endpoint, address)
	if err != nil {
		return authvesting.PeriodicVestingAccount{}, err
	}
	acc, ok := baseAcc.(*authvesting.PeriodicVestingAccount)
	if !ok {
		return authvesting.PeriodicVestingAccount{},
			fmt.Errorf("cannot cast %v to PeriodicVestingAccount", baseAcc)
	}
	return *acc, nil
}

func queryValidator(endpoint, address string) (stakingtypes.Validator, error) {
	var res stakingtypes.QueryValidatorResponse

	body, err := httpGet(fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s", endpoint, address))
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if err := cdc.UnmarshalJSON(body, &res); err != nil {
		return stakingtypes.Validator{}, err
	}
	return res.Validator, nil
}

func queryValidators(endpoint string) (stakingtypes.Validators, error) {
	var res stakingtypes.QueryValidatorsResponse
	body, err := httpGet(fmt.Sprintf("%s/cosmos/staking/v1beta1/validators", endpoint))
	if err != nil {
		return stakingtypes.Validators{}, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if err := cdc.UnmarshalJSON(body, &res); err != nil {
		return stakingtypes.Validators{}, err
	}

	return stakingtypes.Validators{Validators: res.Validators}, nil
}

func queryEvidence(endpoint, hash string) (evidencetypes.QueryEvidenceResponse, error) { //nolint:unused // this is called during e2e tests
	var res evidencetypes.QueryEvidenceResponse
	body, err := httpGet(fmt.Sprintf("%s/cosmos/evidence/v1beta1/evidence/%s", endpoint, hash))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryAllEvidence(endpoint string) (evidencetypes.QueryAllEvidenceResponse, error) {
	var res evidencetypes.QueryAllEvidenceResponse
	body, err := httpGet(fmt.Sprintf("%s/cosmos/evidence/v1beta1/evidence", endpoint))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

//nolint:unused
func queryAllRateLimits(endpoint string) ([]ratelimittypes.RateLimit, error) {
	var res ratelimittypes.QueryAllRateLimitsResponse

	body, err := httpGet(fmt.Sprintf("%s/Stride-Labs/ibc-rate-limiting/ratelimit/ratelimits", endpoint))
	if err != nil {
		return []ratelimittypes.RateLimit{}, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if err := cdc.UnmarshalJSON(body, &res); err != nil {
		return []ratelimittypes.RateLimit{}, err
	}
	return res.RateLimits, nil
}

//nolint:unparam
func queryRateLimit(endpoint, channelID, denom string) (ratelimittypes.QueryRateLimitResponse, error) {
	var res ratelimittypes.QueryRateLimitResponse

	body, err := httpGet(fmt.Sprintf("%s/Stride-Labs/ibc-rate-limiting/ratelimit/ratelimit/%s/by_denom?denom=%s", endpoint, channelID, denom))
	if err != nil {
		return ratelimittypes.QueryRateLimitResponse{}, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if err := cdc.UnmarshalJSON(body, &res); err != nil {
		return ratelimittypes.QueryRateLimitResponse{}, err
	}
	return res, nil
}

func queryRateLimitsByChainID(endpoint, channelID string) ([]ratelimittypes.RateLimit, error) {
	var res ratelimittypes.QueryRateLimitsByChainIdResponse

	body, err := httpGet(fmt.Sprintf("%s/Stride-Labs/ibc-rate-limiting/ratelimit/ratelimits/%s", endpoint, channelID))
	if err != nil {
		return []ratelimittypes.RateLimit{}, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if err := cdc.UnmarshalJSON(body, &res); err != nil {
		return []ratelimittypes.RateLimit{}, err
	}
	return res.RateLimits, nil
}

func queryTradeParams(endpoint string) (tradetypes.QueryParamsResponse, error) {
	body, err := httpGet(fmt.Sprintf("%s/ggezchain/trade/v1beta1/params", endpoint))
	if err != nil {
		return tradetypes.QueryParamsResponse{}, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var params tradetypes.QueryParamsResponse
	if err := cdc.UnmarshalJSON(body, &params); err != nil {
		return tradetypes.QueryParamsResponse{}, err
	}

	return params, nil
}

func queryIBCEscrowAddress(endpoint, channelID string) (string, error) {
	body, err := httpGet(fmt.Sprintf("%s/ibc/apps/transfer/v1/channels/%s/ports/transfer/escrow_address", endpoint, channelID))
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var resp transfertypes.QueryEscrowAddressResponse
	if err := cdc.UnmarshalJSON(body, &resp); err != nil {
		return "", err
	}

	return resp.EscrowAddress, nil
}

func queryICAAccountAddress(endpoint, owner, connectionID string) (string, error) {
	body, err := httpGet(fmt.Sprintf("%s/ibc/apps/interchain_accounts/controller/v1/owners/%s/connections/%s", endpoint, owner, connectionID))
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var icaAccountResp icacontrollertypes.QueryInterchainAccountResponse
	if err := cdc.UnmarshalJSON(body, &icaAccountResp); err != nil {
		return "", err
	}

	return icaAccountResp.Address, nil
}

func querySuperAdmin(endpoint string) (acltypes.QueryGetSuperAdminResponse, error) {
	var res acltypes.QueryGetSuperAdminResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/acl/super_admin", endpoint))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryAllAclAdmin(endpoint string) (acltypes.QueryAllAclAdminResponse, error) {
	var res acltypes.QueryAllAclAdminResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/acl/acl_admin", endpoint))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryAclAdmin(endpoint, address string) (acltypes.QueryGetAclAdminResponse, error) {
	var res acltypes.QueryGetAclAdminResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/acl/acl_admin/%s", endpoint, address))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryAllAclAuthority(endpoint string) (acltypes.QueryAllAclAuthorityResponse, error) {
	var res acltypes.QueryAllAclAuthorityResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/acl/acl_authority", endpoint))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryAclAuthority(endpoint, address string) (acltypes.QueryGetAclAuthorityResponse, error) {
	var res acltypes.QueryGetAclAuthorityResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/acl/acl_authority/%s", endpoint, address))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryAllStoredTrade(endpoint string) (tradetypes.QueryAllStoredTradeResponse, error) {
	var res tradetypes.QueryAllStoredTradeResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/trade/stored_trade", endpoint))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryStoredTrade(endpoint, tradeIndex string) (tradetypes.QueryGetStoredTradeResponse, error) {
	var res tradetypes.QueryGetStoredTradeResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/trade/stored_trade/%s", endpoint, tradeIndex))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryAllStoredTempTrade(endpoint string) (tradetypes.QueryAllStoredTempTradeResponse, error) {
	var res tradetypes.QueryAllStoredTempTradeResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/trade/stored_temp_trade", endpoint))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryStoredTempTrade(endpoint, tradeIndex string) (tradetypes.QueryGetStoredTempTradeResponse, error) {
	var res tradetypes.QueryGetStoredTempTradeResponse

	body, err := httpGet(fmt.Sprintf("%s/GGEZLabs/ggezchain/trade/stored_temp_trade/%s", endpoint, tradeIndex))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}
