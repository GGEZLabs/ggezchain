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
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
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
			valid: true,
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
			valid: false,
		},
		{
			desc: "duplicated storedTempTrade",
			genState: &types.GenesisState{
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
						Amount:     &sdk.Coin{Denom: "invalid_denom", Amount: math.NewInt(0)},
					},
				},
			},
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
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
			valid: false,
		},
		{
			desc: "invalid next_id",
			genState: &types.GenesisState{
				TradeIndex: types.TradeIndex{
					NextId: 0,
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
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
