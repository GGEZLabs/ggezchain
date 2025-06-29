package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccessDefinitionsJson(t *testing.T) {
	t.Run("nil AccessDefinitions", func(t *testing.T) {
		acl := AclAuthority{
			AccessDefinitions: nil,
		}
		result := acl.AccessDefinitionsJson()
		require.Equal(t, "[]", result)
	})

	t.Run("empty AccessDefinitions", func(t *testing.T) {
		acl := AclAuthority{
			AccessDefinitions: []*AccessDefinition{},
		}
		result := acl.AccessDefinitionsJson()
		require.Equal(t, "[]", result)
	})

	t.Run("single AccessDefinition", func(t *testing.T) {
		acl := AclAuthority{
			AccessDefinitions: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
			},
		}
		expected, _ := json.Marshal(acl.AccessDefinitions)
		result := acl.AccessDefinitionsJson()
		require.Equal(t, string(expected), result)
	})

	t.Run("multiple AccessDefinitions", func(t *testing.T) {
		acl := AclAuthority{
			AccessDefinitions: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
				{Module: "module2", IsMaker: false, IsChecker: true},
			},
		}
		expected, _ := json.Marshal(acl.AccessDefinitions)
		result := acl.AccessDefinitionsJson()
		require.Equal(t, string(expected), result)
	})
}
