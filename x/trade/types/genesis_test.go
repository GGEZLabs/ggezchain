package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "{}",
					},
				},
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:44:00Z",
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "{}",
					},
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "{}",
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
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:44:00Z",
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
				StoredTradeList: []types.StoredTrade{
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeUnspecified,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_type must be buy or sell",
		},
		{
			desc: "invalid amount",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(0)},
					},
				},
			},
			expErr:    true,
			expErrMsg: "amount must be more than 0",
		},
		{
			desc: "invalid denom",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTradeList: []types.StoredTrade{
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
			desc: "invalid price",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:      "0",
					},
				},
			},
			expErr:    true,
			expErrMsg: "price must be more than 0",
		},
		{
			desc: "invalid receiver_address",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid receiver_address",
		},
		{
			desc: "invalid status",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusNil,
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           "invalid_address",
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid checker address",
		},
		{
			desc: "invalid create_date",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-1108:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid create_date format",
		},
		{
			desc: "invalid update_date",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-11T08:44:00Z",
						UpdateDate:      "2023-05-11T08",
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-11T08:44:00Z",
						UpdateDate:      "2023-05-11T08:44:00Z",
						ProcessDate:     "2023-05-0T08:44:00Z",
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-11T08:44:00Z",
						UpdateDate:      "2023-05-11T08:44:00Z",
						ProcessDate:     "2023-05-11T08:44:00Z",
						TradeData:       "invalid_trade_data",
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid banking_system_data JSON format",
		},
		{
			desc: "invalid trade_index (stored temp trade)",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_index must be more than 0",
		},
		{
			desc: "invalid temp_trade_index (stored temp trade)",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "temp_trade_index must be more than 0",
		},
		{
			desc: "duplicate temp_trade_index (stored temp trade)",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex:     2,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-12T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated temp_trade_index",
		},
		{
			desc: "invalid create_date (stored temp trade)",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 2,
				},
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:4",
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
	tests := []struct {
		desc      string
		genState  *types.GenesisState
		expErr    bool
		expErrMsg string
	}{
		{
			desc: "valid storedTrade list",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "{}",
					},
				},
			},
			expErr: false,
		},
		{
			desc: "duplicated storedTrade",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "{}",
					},
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "{}",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated index for storedTrade",
		},
		{
			desc: "invalid trade_index",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
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
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeUnspecified,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_type must be buy or sell",
		},
		{
			desc: "invalid amount",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(0)},
					},
				},
			},
			expErr:    true,
			expErrMsg: "amount must be more than 0",
		},
		{
			desc: "invalid denom",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
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
			desc: "invalid price",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex: 1,
						TradeType:  types.TradeTypeBuy,
						Amount:     &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:      "0",
					},
				},
			},
			expErr:    true,
			expErrMsg: "price must be more than 0",
		},
		{
			desc: "invalid receiver_address",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid receiver_address",
		},
		{
			desc: "invalid status",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusNil,
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid status",
		},
		{
			desc: "invalid maker address",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid maker address",
		},
		{
			desc: "invalid checker address",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         "invalid_address",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid checker address",
		},
		{
			desc: "invalid create_date",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-1108:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid create_date format",
		},
		{
			desc: "invalid update_date",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-11T08:44:00Z",
						UpdateDate:      "2023-05-11T08",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid update_date format",
		},
		{
			desc: "invalid process_date",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-11T08:44:00Z",
						UpdateDate:      "2023-05-11T08:44:00Z",
						ProcessDate:     "2023-05-0T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid process_date format",
		},
		{
			desc: "invalid trade_data",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:      1,
						TradeType:       types.TradeTypeBuy,
						Amount:          &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:           "0.01",
						ReceiverAddress: sample.AccAddress(),
						Status:          types.StatusPending,
						Maker:           sample.AccAddress(),
						Checker:         sample.AccAddress(),
						CreateDate:      "2023-05-11T08:44:00Z",
						UpdateDate:      "2023-05-11T08:44:00Z",
						ProcessDate:     "2023-05-11T08:44:00Z",
						TradeData:       "invalid_trade_data",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid trade_data",
		},
		{
			desc: "invalid banking_system_data",
			genState: &types.GenesisState{
				StoredTradeList: []types.StoredTrade{
					{
						TradeIndex:        1,
						TradeType:         types.TradeTypeBuy,
						Amount:            &sdk.Coin{Denom: types.DefaultCoinDenom, Amount: math.NewInt(100000)},
						Price:             "0.01",
						ReceiverAddress:   sample.AccAddress(),
						Status:            types.StatusPending,
						Maker:             sample.AccAddress(),
						Checker:           sample.AccAddress(),
						CreateDate:        "2023-05-11T08:44:00Z",
						UpdateDate:        "2023-05-11T08:44:00Z",
						ProcessDate:       "2023-05-11T08:44:00Z",
						TradeData:         types.GetSampleTradeData(),
						BankingSystemData: "",
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid banking_system_data JSON format",
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
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex:     2,
						TempTradeIndex: 2,
						CreateDate:     "2023-05-11T10:44:00Z",
					},
				},
			},
			expErr: false,
		},
		{
			desc: "invalid trade_index",
			genState: &types.GenesisState{
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "trade_index must be more than 0",
		},
		{
			desc: "invalid temp_trade_index",
			genState: &types.GenesisState{
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 0,
					},
				},
			},
			expErr:    true,
			expErrMsg: "temp_trade_index must be more than 0",
		},
		{
			desc: "duplicate temp_trade_index",
			genState: &types.GenesisState{
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:44:00Z",
					},
					{
						TradeIndex:     2,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-12T08:44:00Z",
					},
				},
			},
			expErr:    true,
			expErrMsg: "duplicated temp_trade_index",
		},
		{
			desc: "invalid create_date",
			genState: &types.GenesisState{
				StoredTempTradeList: []types.StoredTempTrade{
					{
						TradeIndex:     1,
						TempTradeIndex: 1,
						CreateDate:     "2023-05-11T08:4",
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
			StoredTradeList:     []types.StoredTrade{},
			TradeIndex:          types.TradeIndex{NextId: 1},
			Params:              types.Params{},
			StoredTempTradeList: []types.StoredTempTrade{},
		},
		types.DefaultGenesis())
}
