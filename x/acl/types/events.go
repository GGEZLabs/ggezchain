package types

// acl module event types
const (
	EventTypeAddAuthority    = "add_authority"
	EventTypeDeleteAuthority = "delete_authority"
	EventTypeUpdateAuthority = "update_authority"
	EventTypeInit            = "init"
	EventTypeAddAdmin        = "add_admin"
	EventTypeDeleteAdmin     = "delete_admin"

	AttributeKeyAuthorityAddress  = "authority_address"
	AttributeKeyName              = "name"
	AttributeKeyAccessDefinitions = "access_definitions"
	AttributeKeyAdmins            = "admins"
)
