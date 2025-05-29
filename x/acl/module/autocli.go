package acl

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	modulev1 "github.com/GGEZLabs/ggezchain/api/ggezchain/acl"
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
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "AclAuthorityAll",
					Use:       "list-acl-authority",
					Short:     "List all aclAuthority",
				},
				{
					RpcMethod:      "AclAuthority",
					Use:            "show-acl-authority [address]",
					Short:          "Shows a aclAuthority",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "AclAdminAll",
					Use:       "list-acl-admin",
					Short:     "List all aclAdmin",
				},
				{
					RpcMethod:      "AclAdmin",
					Use:            "show-acl-admin [address]",
					Short:          "Shows a aclAdmin",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "SuperAdmin",
					Use:       "super-admin",
					Short:     "show super_admin",
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
					RpcMethod:      "AddAuthority",
					Use:            "add-authority [auth-address] [name] [access-definitions]",
					Short:          "Add a new authority with specific access definition. Must have admin authority to do so.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auth_address"}, {ProtoField: "name"}, {ProtoField: "access_definitions"}},
				},
				{
					RpcMethod:      "DeleteAuthority",
					Use:            "delete-authority [auth-address]",
					Short:          "Delete an existing authority. Must have admin authority to do so.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auth_address"}},
				},
				{
					RpcMethod:      "UpdateAuthority",
					Use:            "update-authority [auth-address]",
					Short:          "Modify the name or access definition of an existing authority. Must have admin authority to do so.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auth_address"}},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"new_name": {
							Name:         "new-name",
							Usage:        "Set a new name for the authority.",
							DefaultValue: "",
						},
						"overwrite_access_definitions": {
							Name:         "overwrite-access-definitions",
							Usage:        "Overwrite the entire access definition list with this JSON string. Ignores other access definition flags.",
							DefaultValue: "",
						},
						"add_access_definitions": {
							Name:         "add-access-definitions",
							Usage:        "Add one or more new access definition.",
							DefaultValue: "",
						},
						"update_access_definition": {
							Name:         "update-access-definition",
							Usage:        "Update access definition values for an existing module. (matched by module name)",
							DefaultValue: "",
						},
						"delete_access_definitions": {
							Name:         "delete-access-definitions",
							Usage:        "Delete one or more specific access definition (by module name).",
							DefaultValue: "",
						},
						"clear_all_access_definitions": {
							Name:         "clear-all-access-definitions",
							Usage:        "Clear all access definition. Default is false.",
							DefaultValue: "false",
						},
					},
					Example: `Overwrite the entire access definition list with this JSON string. Ignores other access definition flags:
ggezchaind tx acl update-authority ggezauthaddress... --add-access-definitions '[{"module":"module1","is_maker":true,"is_checker":false}]' --from ggezaddress...

Add one or more new access definition:
ggezchaind tx acl update-authority ggezauthaddress... --add-access-definitions '[{"module":"module2","is_maker":true,"is_checker":true}]' --from ggezaddress...

Update access definition values for an existing module (by module name):
ggezchaind tx acl update-authority ggezauthaddress... --update-access-definition '{"module":"module2","is_maker":false,"is_checker":true}' --from ggezaddress...

Delete one or more specific access definition (by module name):
ggezchaind tx acl update-authority ggezauthaddress... --delete-access-definitions 'module2,module1' --from ggezaddress...

Clear all access definition. (Default is false)
ggezchaind tx acl update-authority ggezauthaddress... --clear-all-access-definitions --from ggezaddress...

`,
				},
				{
					RpcMethod:      "Init",
					Use:            "init [admins]",
					Short:          "Initializes the AclAdmin by one or more admin. Can only be called once.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admins"}},
				},
				{
					RpcMethod:      "AddAdmin",
					Use:            "add-admin [admins]",
					Short:          "Add one or more aclAdmin. Must have admin authority to do so.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admins"}},
				},
				{
					RpcMethod:      "DeleteAdmin",
					Use:            "delete-admin [admins]",
					Short:          "Delete one or more aclAdmin. Must have admin authority to do so.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admins"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
