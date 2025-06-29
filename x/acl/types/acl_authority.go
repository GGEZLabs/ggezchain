package types

import "encoding/json"

// AccessDefinitionsJson returns the JSON-encoded string of the AccessDefinitions field.
// If the field is nil or marshalling fails, it returns an empty JSON array "[]"
func (acl *AclAuthority) AccessDefinitionsJson() string {
	if acl.AccessDefinitions == nil {
		return "[]"
	}

	data, err := json.Marshal(acl.AccessDefinitions)
	if err != nil {
		return "[]"
	}

	return string(data)
}
