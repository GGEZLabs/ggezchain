package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	td := types.GetSampleTradeData(types.TradeTypeBuy)
	cmpj := types.GetSampleCoinMintingPriceJson()
	erj := types.GetSampleExchangeRateJson()
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			expErr:   false,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: cmpj,
						ExchangeRateJson:     erj,
					},
				},
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:44:00Z",
					},
				},
			},
			expErr: false,
		},
		{
			desc: "duplicated storedTrade",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 1,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: cmpj,
						ExchangeRateJson:     erj,
					},
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: cmpj,
						ExchangeRateJson:     erj,
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for storedTrade",
		},
		{
			desc: "duplicated storedTempTrade",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 1,
				},
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for storedTempTrade",
		},
		{
			desc: "invalid trade_index",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_index must be more than 0",
		},
		{
			desc: "invalid trade_type",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeUnspecified,
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid trade_type",
		},
		{
			desc: "invalid amount",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(0)},
					},
				},
			},
			expErr:    true,
			expErrMsg: "zero amount not allowed",
		},
		{
			desc: "invalid denom",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: "invalid_denom", Amount: math.NewInt(10)},
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid denom",
		},
		{
			desc: "invalid receiver_address",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusPending,
						ReceiverAddress:     "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid receiver_address",
		},
		{
			desc: "set amount with trade type split",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeSplit,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusPending,
					},
				},
			},
			expErr:    true,
			expErrMsg: "amount must not be set for trade type: TRADE_TYPE_SPLIT",
		},
		{
			desc: "set receiver address with trade type split",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeSplit,
						Amount:              &sdk.Coin{Denom: "", Amount: math.NewInt(0)},
						ReceiverAddress:     sample.AccAddress(),
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusPending,
					},
				},
			},
			expErr:    true,
			expErrMsg: "receiver_address must not be set for trade type TRADE_TYPE_SPLIT",
		},
		{
			desc: "invalid price",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						ReceiverAddress:     sample.AccAddress(),
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0",
					},
				},
			},
			expErr:    true,
			expErrMsg: "price must be more than 0",
		},
		{
			desc: "invalid status",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						ReceiverAddress:     sample.AccAddress(),
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusNil,
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid status",
		},
		{
			desc: "invalid maker address",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusPending,
						Maker:               "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid maker address",
		},
		{
			desc: "invalid checker address",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusRejected,
						Maker:               sample.AccAddress(),
						Checker:             "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid checker address",
		},
		{
			desc: "set checker address with status pending",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusPending,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
					},
				},
			},
			expErr:    true,
			expErrMsg: "checker must not be set for trade status TRADE_STATUS_PENDING",
		},
		{
			desc: "invalid create_date",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-1108:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid create_date format",
		},
		{
			desc: "invalid tx_date",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-1108:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid tx_date format",
		},
		{
			desc: "invalid update_date",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid update_date format",
		},
		{
			desc: "invalid process_date",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08:44:00Z",
						ProcessDate:         "2023-05-0T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid process_date format",
		},
		{
			desc: "invalid trade_data",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08:44:00Z",
						ProcessDate:         "2023-05-11T08:44:00Z",
						TradeData:           "invalid_trade_data",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid trade_data",
		},
		{
			desc: "invalid banking_system_data",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08:44:00Z",
						ProcessDate:         "2023-05-11T08:44:00Z",
						TradeData:           td,
						BankingSystemData:   "",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid banking_system_data json format",
		},
		{
			desc: "invalid coin_minting_price_json",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: "invalid_coin_minting_price_json",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid coin_minting_price_json",
		},
		{
			desc: "invalid exchange_rate_json",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
						ExchangeRateJson:     "invalid_exchange_rate_json",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid exchange_rate_json",
		},
		{
			desc: "invalid trade_index (stored temp trade)",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_index must be more than 0",
		},
		{
			desc: "duplicate trade_index (stored temp trade)",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex: 1,
						TxDate:     "2023-05-12T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for storedTempTrade",
		},
		{
			desc: "invalid create_date (stored temp trade)",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:4",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid create_date format",
		},
		{
			desc: "invalid next_id",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 0,
				},
			},
			expErr:    true,
			expErrMsg: "next_id must be more than 0",
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenesisState_ValidateStoredTrade(t *testing.T) {
	td := types.GetSampleTradeData(types.TradeTypeBuy)
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
	}{
		{
			desc: "valid storedTrade list",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
						ExchangeRateJson:     types.GetSampleExchangeRateJson(),
					},
				},
			},
			expErr: false,
		},
		{
			desc: "duplicated storedTrade",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
						ExchangeRateJson:     types.GetSampleExchangeRateJson(),
					},
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusPending,
						Maker:                sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
						ExchangeRateJson:     types.GetSampleExchangeRateJson(),
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for storedTrade",
		},
		{
			desc: "invalid trade_index",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_index must be more than 0",
		},
		{
			desc: "invalid trade_type",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeUnspecified,
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid trade_type",
		},
		{
			desc: "invalid amount",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(0)},
					},
				},
			},
			expErr:    true,
			expErrMsg: "zero amount not allowed",
		},
		{
			desc: "invalid denom",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: "invalid_denom", Amount: math.NewInt(10)},
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid denom",
		},
		{
			desc: "invalid receiver_address",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusPending,
						ReceiverAddress:     "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid receiver_address",
		},
		{
			desc: "set amount with trade type split",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeSplit,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusPending,
					},
				},
			},
			expErr:    true,
			expErrMsg: "amount must not be set for trade type: TRADE_TYPE_SPLIT",
		},
		{
			desc: "set receiver address with trade type split",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeSplit,
						Amount:              &sdk.Coin{Denom: "", Amount: math.NewInt(0)},
						ReceiverAddress:     sample.AccAddress(),
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusPending,
					},
				},
			},
			expErr:    true,
			expErrMsg: "receiver_address must not be set for trade type TRADE_TYPE_SPLIT",
		},
		{
			desc: "invalid price",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						ReceiverAddress:     sample.AccAddress(),
						CoinMintingPriceUsd: "0",
					},
				},
			},
			expErr:    true,
			expErrMsg: "price must be more than 0",
		},
		{
			desc: "invalid status",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						ReceiverAddress:     sample.AccAddress(),
						CoinMintingPriceUsd: "0.01",
						Status:              types.StatusNil,
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid status",
		},
		{
			desc: "invalid maker address",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusPending,
						Maker:               "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid maker address",
		},
		{
			desc: "invalid checker address",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusRejected,
						Maker:               sample.AccAddress(),
						Checker:             "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid checker address",
		},
		{
			desc: "set checker address with status pending",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusPending,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
					},
				},
			},
			expErr:    true,
			expErrMsg: "checker must not be set for trade status TRADE_STATUS_PENDING",
		},
		{
			desc: "invalid create_date",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-1108:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid create_date format",
		},
		{
			desc: "invalid tx_date",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-1108:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid tx_date format",
		},
		{
			desc: "invalid update_date",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid update_date format",
		},
		{
			desc: "invalid process_date",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08:44:00Z",
						ProcessDate:         "2023-05-0T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid process_date format",
		},
		{
			desc: "invalid trade_data",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusRejected,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08:44:00Z",
						ProcessDate:         "2023-05-11T08:44:00Z",
						TradeData:           "invalid_trade_data",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid trade_data",
		},
		{
			desc: "invalid banking_system_data",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:          1,
						TradeType:           types.TradeTypeBuy,
						Amount:              &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd: "0.01",
						ReceiverAddress:     sample.AccAddress(),
						Status:              types.StatusProcessed,
						Maker:               sample.AccAddress(),
						Checker:             sample.AccAddress(),
						CreateDate:          "2023-05-11T08:44:00Z",
						TxDate:              "2023-05-11T08:44:00Z",
						UpdateDate:          "2023-05-11T08:44:00Z",
						ProcessDate:         "2023-05-11T08:44:00Z",
						TradeData:           td,
						BankingSystemData:   "",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid banking_system_data json format",
		},
		{
			desc: "invalid coin_minting_price_json",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: "invalid_coin_minting_price_json",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid coin_minting_price_json",
		},
		{
			desc: "invalid exchange_rate_json",
			genState: &types.GenesisState{
				StoredTrades: []types.StoredTrade{
					{
						TradeIndex:           1,
						TradeType:            types.TradeTypeBuy,
						Amount:               &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
						CoinMintingPriceUsd:  "0.01",
						ReceiverAddress:      sample.AccAddress(),
						Status:               types.StatusProcessed,
						Maker:                sample.AccAddress(),
						Checker:              sample.AccAddress(),
						CreateDate:           "2023-05-11T08:44:00Z",
						TxDate:               "2023-05-11T08:44:00Z",
						UpdateDate:           "2023-05-11T08:44:00Z",
						ProcessDate:          "2023-05-11T08:44:00Z",
						TradeData:            td,
						BankingSystemData:    "{}",
						CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
						ExchangeRateJson:     "invalid_exchange_rate_json",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid exchange_rate_json",
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.ValidateStoredTrade()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenesisState_ValidateStoredTempTrade(t *testing.T) {
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
	}{
		{
			desc: "valid storedTempTrade list",
			genState: &types.GenesisState{
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex: 2,
						TxDate:     "2023-05-11T10:44:00Z",
					},
				},
			},
			expErr: false,
		},
		{
			desc: "invalid trade_index",
			genState: &types.GenesisState{
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_index must be more than 0",
		},
		{
			desc: "duplicate trade_index",
			genState: &types.GenesisState{
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex: 1,
						TxDate:     "2023-05-12T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for storedTempTrade",
		},
		{
			desc: "invalid create_date",
			genState: &types.GenesisState{
				StoredTempTrades: []types.StoredTempTrade{
					{
						TradeIndex: 1,
						TxDate:     "2023-05-11T08:4",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid create_date format",
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.ValidateStoredTempTrade()
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDefaultGenesisState_ExpectedInitialNextId(t *testing.T) {
	require.EqualValues(t,
		&types.GenesisState{
			StoredTrades:     []types.StoredTrade{},
			TradeIndex:       types.TradeIndex{NextId: 1},
			Params:           types.Params{},
			StoredTempTrades: []types.StoredTempTrade{},
		},
		types.DefaultGenesis())
}
