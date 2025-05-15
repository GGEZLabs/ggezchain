package types

// acl module event types
const (
	EventTypeAddAuthority    = "add_authority"
	EventTypeDeleteAuthority = "delete_authority"
	EventTypeUpdateAuthority = "update_authority"
	EventTypeInitAclAdmin    = "init_acl_admin"
	EventTypeAddAclAdmin     = "add_acl_admin"
	EventTypeDeleteAclAdmin  = "delete_acl_admin"

	AttributeKeyAuthorityAddress  = "authority_address"
	AttributeKeyName              = "name"
	AttributeKeyAccessDefinitions = "access_definitions"
	AttributeKeyAdmins            = "admins"
)
