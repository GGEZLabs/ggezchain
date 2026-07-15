package trade

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the parameters of the module",
				},
				{
					RpcMethod: "GetTradeIndex",
					Use:       "trade-index",
					Short:     "Query a trade-index",
				},
				{
					RpcMethod: "ListStoredTrade",
					Use:       "stored-trades",
					Short:     "Query all stored-trades",
				},
				{
					RpcMethod:      "GetStoredTrade",
					Use:            "stored-trade [trade-index]",
					Short:          "Query a stored-trade by trade index",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}},
				},
				{
					RpcMethod: "ListStoredTempTrade",
					Use:       "stored-temp-trades",
					Short:     "Query all stored-temp-trades",
				},
				{
					RpcMethod:      "GetStoredTempTrade",
					Use:            "stored-temp-trade [trade-index]",
					Short:          "Query a stored-temp-trade by trade index",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "CreateTrade",
					Use:       "create-trade [trade-data] [banking-system-data] [coin-minting-price-json] [exchange-rate-json] [receiver-address]",
					Short:     "Create the StoredTrade. Must have authority to do so.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "trade_data"},
						{ProtoField: "banking_system_data"},
						{ProtoField: "coin_minting_price_json"},
						{ProtoField: "exchange_rate_json"},
						{ProtoField: "receiver_address", Optional: true},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"create_date": {
							Name:         "create-date",
							Usage:        "Set a create date. Default is current date",
							DefaultValue: "",
						},
					},
				},
				{
					RpcMethod:      "ProcessTrade",
					Use:            "process-trade [trade-index] [process-type]",
					Short:          "Process the StoredTrade. Must have authority to do so.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}, {ProtoField: "process_type"}},
				},
			},
		},
	}
}
