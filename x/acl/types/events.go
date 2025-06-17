package types

// acl module event types
const (
	EventTypeAddAuthority     = "add_authority"
	EventTypeDeleteAuthority  = "delete_authority"
	EventTypeUpdateAuthority  = "update_authority"
	EventTypeInit             = "init"
	EventTypeUpdateSuperAdmin = "update_super_admin"
	EventTypeAddAdmin         = "add_admin"
	EventTypeDeleteAdmin      = "delete_admin"

	AttributeKeySuperAdmin        = "super_admin"
	AttributeKeyAuthorityAddress  = "authority_address"
	AttributeKeyName              = "name"
	AttributeKeyAccessDefinitions = "access_definitions"
	AttributeKeyAdmins            = "admins"
)
