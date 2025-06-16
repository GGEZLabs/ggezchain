package trade

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	modulev1 "github.com/GGEZLabs/ggezchain/api/ggezchain/trade"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the parameters of the module",
				},
				{
					RpcMethod: "TradeIndex",
					Use:       "trade-index",
					Short:     "Query a trade-index",
				},
				{
					RpcMethod: "StoredTradeAll",
					Use:       "stored-trades",
					Short:     "Query all stored-trades",
				},
				{
					RpcMethod:      "StoredTrade",
					Use:            "stored-trade [trade-index]",
					Short:          "Query a stored-trade by trade index",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}},
				},
				{
					RpcMethod: "StoredTempTradeAll",
					Use:       "stored-temp-trades",
					Short:     "Query all stored-temp-trades",
				},
				{
					RpcMethod:      "StoredTempTrade",
					Use:            "stored-temp-trade [trade-index]",
					Short:          "Query a stored-temp-trade by trade index",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "CreateTrade",
					Use:       "create-trade [receiver-address] [trade-data] [banking-system-data] [coin-minting-price-json] [exchange-rate-json]",
					Short:     "Send a create-trade tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "receiver_address"},
						{ProtoField: "trade_data"},
						{ProtoField: "banking_system_data"},
						{ProtoField: "coin_minting_price_json"},
						{ProtoField: "exchange_rate_json"},
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
					Short:          "Send a process-trade tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}, {ProtoField: "process_type"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
