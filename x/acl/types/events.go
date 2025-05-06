package types

// Trade module event types
const (
	EventTypeAddAuthority    = "add_authority"
	EventTypeDeleteAuthority = "delete_authority"
	EventTypeUpdateAuthority = "update_authority"

	AttributeKeyAuthorityAddress  = "authority_address"
	AttributeKeyName              = "name"
	AttributeKeyAccessDefinitions = "access_definitions"
)
