package acl

import (
	"math/rand"

	aclsimulation "github.com/GGEZLabs/ggezchain/v2/x/acl/simulation"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
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
		opWeightMsgInit          = "op_weight_msg_init_acl_admin"
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
		opWeightMsgUpdateSuperAdmin          = "op_weight_msg_update_super_admin"
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
		opWeightMsgAddAdmin          = "op_weight_msg_add_acl_admin"
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
		opWeightMsgDeleteAdmin          = "op_weight_msg_delete_acl_admin"
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
		opWeightMsgAddAuthority          = "op_weight_msg_add_authority"
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
		opWeightMsgUpdateAuthority          = "op_weight_msg_update_authority"
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
	const (
		opWeightMsgDeleteAuthority          = "op_weight_msg_delete_authority"
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

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
//
// Unlike the tx-operation simulators above (which submit and deliver a full
// transaction), each MsgSimulatorFn here only needs to construct a plausible
// concrete sdk.Msg for x/gov's proposal simulation to wrap into a governance
// proposal; returning nil tells the gov harness to skip that proposal message
// for this round (used below when there's no existing admin/authority to
// target yet).
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	const (
		opWeightMsgInit          = "op_weight_msg_init_acl_admin"
		defaultWeightMsgInit int = 100

		opWeightMsgUpdateSuperAdmin          = "op_weight_msg_update_super_admin"
		defaultWeightMsgUpdateSuperAdmin int = 100

		opWeightMsgAddAdmin          = "op_weight_msg_add_acl_admin"
		defaultWeightMsgAddAdmin int = 100

		opWeightMsgDeleteAdmin          = "op_weight_msg_delete_acl_admin"
		defaultWeightMsgDeleteAdmin int = 100

		opWeightMsgAddAuthority          = "op_weight_msg_add_authority"
		defaultWeightMsgAddAuthority int = 100

		opWeightMsgUpdateAuthority          = "op_weight_msg_update_authority"
		defaultWeightMsgUpdateAuthority int = 100

		opWeightMsgDeleteAuthority          = "op_weight_msg_delete_authority"
		defaultWeightMsgDeleteAuthority int = 100
	)

	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgInit,
			defaultWeightMsgInit,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgInit(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateSuperAdmin,
			defaultWeightMsgUpdateSuperAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgUpdateSuperAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddAdmin,
			defaultWeightMsgAddAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgAddAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteAdmin,
			defaultWeightMsgDeleteAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgDeleteAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddAuthority,
			defaultWeightMsgAddAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgAddAuthority(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateAuthority,
			defaultWeightMsgUpdateAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgUpdateAuthority(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteAuthority,
			defaultWeightMsgDeleteAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgDeleteAuthority(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
	}
}
