package trade

import (
	"math/rand"

	tradesimulation "github.com/GGEZLabs/ggezchain/v2/x/trade/simulation"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tradeGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		TradeIndex: types.TradeIndex{
			NextId: uint64(simState.Rand.Intn(99) + 1),
		},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tradeGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateTrade          = "op_weight_msg_create_trade"
		defaultWeightMsgCreateTrade int = 100
	)

	var weightMsgCreateTrade int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTrade, &weightMsgCreateTrade, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTrade = defaultWeightMsgCreateTrade
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTrade,
		tradesimulation.SimulateMsgCreateTrade(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgProcessTrade          = "op_weight_msg_process_trade"
		defaultWeightMsgProcessTrade int = 100
	)

	var weightMsgProcessTrade int
	simState.AppParams.GetOrGenerate(opWeightMsgProcessTrade, &weightMsgProcessTrade, nil,
		func(_ *rand.Rand) {
			weightMsgProcessTrade = defaultWeightMsgProcessTrade
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgProcessTrade,
		tradesimulation.SimulateMsgProcessTrade(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
//
// Unlike the tx-operation simulators above (which submit and deliver a full
// transaction), each MsgSimulatorFn here only needs to construct a plausible
// concrete sdk.Msg for x/gov's proposal simulation to wrap into a governance
// proposal; returning nil tells the gov harness to skip that proposal message
// for this round (used below for MsgProcessTrade when there's no pending trade
// yet to target).
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	const (
		opWeightMsgCreateTrade          = "op_weight_msg_create_trade"
		defaultWeightMsgCreateTrade int = 100

		opWeightMsgProcessTrade          = "op_weight_msg_process_trade"
		defaultWeightMsgProcessTrade int = 100
	)

	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateTrade,
			defaultWeightMsgCreateTrade,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradesimulation.SimulateMsgCreateTrade(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgProcessTrade,
			defaultWeightMsgProcessTrade,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradesimulation.SimulateMsgProcessTrade(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
	}
}
