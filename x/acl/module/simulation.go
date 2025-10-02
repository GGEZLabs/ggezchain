package acl

import (
	"math/rand"

	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	aclsimulation "github.com/GGEZLabs/ggezchain/v2/x/acl/simulation"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = aclsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgAddAuthority = "op_weight_msg_add_authority"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddAuthority int = 100

	opWeightMsgDeleteAuthority = "op_weight_msg_delete_authority"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAuthority int = 100

	opWeightMsgUpdateAuthority = "op_weight_msg_update_authority"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateAuthority int = 100

	opWeightMsgInit = "op_weight_msg_init_acl_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgInit int = 100

	opWeightMsgAddAdmin = "op_weight_msg_add_acl_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddAdmin int = 100

	opWeightMsgDeleteAdmin = "op_weight_msg_delete_acl_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAdmin int = 100

	opWeightMsgUpdateSuperAdmin = "op_weight_msg_update_super_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateSuperAdmin int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	aclGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&aclGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAddAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgAddAuthority, &weightMsgAddAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgAddAuthority = defaultWeightMsgAddAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddAuthority,
		aclsimulation.SimulateMsgAddAuthority(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgDeleteAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAuthority, &weightMsgDeleteAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAuthority = defaultWeightMsgDeleteAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAuthority,
		aclsimulation.SimulateMsgDeleteAuthority(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgUpdateAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAuthority, &weightMsgUpdateAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAuthority = defaultWeightMsgUpdateAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAuthority,
		aclsimulation.SimulateMsgUpdateAuthority(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgInit int
	simState.AppParams.GetOrGenerate(opWeightMsgInit, &weightMsgInit, nil,
		func(_ *rand.Rand) {
			weightMsgInit = defaultWeightMsgInit
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgInit,
		aclsimulation.SimulateMsgInit(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgAddAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgAddAdmin, &weightMsgAddAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgAddAdmin = defaultWeightMsgAddAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddAdmin,
		aclsimulation.SimulateMsgAddAdmin(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgDeleteAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAdmin, &weightMsgDeleteAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAdmin = defaultWeightMsgDeleteAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAdmin,
		aclsimulation.SimulateMsgDeleteAdmin(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgUpdateSuperAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateSuperAdmin, &weightMsgUpdateSuperAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateSuperAdmin = defaultWeightMsgUpdateSuperAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateSuperAdmin,
		aclsimulation.SimulateMsgUpdateSuperAdmin(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddAuthority,
			defaultWeightMsgAddAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgAddAuthority(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteAuthority,
			defaultWeightMsgDeleteAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgDeleteAuthority(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateAuthority,
			defaultWeightMsgUpdateAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgUpdateAuthority(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgInit,
			defaultWeightMsgInit,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgInit(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddAdmin,
			defaultWeightMsgAddAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgAddAdmin(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteAdmin,
			defaultWeightMsgDeleteAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgDeleteAdmin(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateSuperAdmin,
			defaultWeightMsgUpdateSuperAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgUpdateSuperAdmin(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
