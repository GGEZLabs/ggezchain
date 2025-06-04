package acl

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	aclsimulation "github.com/GGEZLabs/ramichain/x/acl/simulation"
	"github.com/GGEZLabs/ramichain/x/acl/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	aclGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&aclGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgInit          = "op_weight_msg_acl"
		defaultWeightMsgInit int = 100
	)

	var weightMsgInit int
	simState.AppParams.GetOrGenerate(opWeightMsgInit, &weightMsgInit, nil,
		func(_ *rand.Rand) {
			weightMsgInit = defaultWeightMsgInit
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgInit,
		aclsimulation.SimulateMsgInit(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgAddAdmin          = "op_weight_msg_acl"
		defaultWeightMsgAddAdmin int = 100
	)

	var weightMsgAddAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgAddAdmin, &weightMsgAddAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgAddAdmin = defaultWeightMsgAddAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddAdmin,
		aclsimulation.SimulateMsgAddAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteAdmin          = "op_weight_msg_acl"
		defaultWeightMsgDeleteAdmin int = 100
	)

	var weightMsgDeleteAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAdmin, &weightMsgDeleteAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAdmin = defaultWeightMsgDeleteAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAdmin,
		aclsimulation.SimulateMsgDeleteAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateSuperAdmin          = "op_weight_msg_acl"
		defaultWeightMsgUpdateSuperAdmin int = 100
	)

	var weightMsgUpdateSuperAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateSuperAdmin, &weightMsgUpdateSuperAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateSuperAdmin = defaultWeightMsgUpdateSuperAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateSuperAdmin,
		aclsimulation.SimulateMsgUpdateSuperAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgAddAuthority          = "op_weight_msg_acl"
		defaultWeightMsgAddAuthority int = 100
	)

	var weightMsgAddAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgAddAuthority, &weightMsgAddAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgAddAuthority = defaultWeightMsgAddAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddAuthority,
		aclsimulation.SimulateMsgAddAuthority(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteAuthority          = "op_weight_msg_acl"
		defaultWeightMsgDeleteAuthority int = 100
	)

	var weightMsgDeleteAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAuthority, &weightMsgDeleteAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAuthority = defaultWeightMsgDeleteAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAuthority,
		aclsimulation.SimulateMsgDeleteAuthority(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateAuthority          = "op_weight_msg_acl"
		defaultWeightMsgUpdateAuthority int = 100
	)

	var weightMsgUpdateAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAuthority, &weightMsgUpdateAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAuthority = defaultWeightMsgUpdateAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAuthority,
		aclsimulation.SimulateMsgUpdateAuthority(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
