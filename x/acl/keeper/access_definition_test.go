package keeper_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/ggezchain/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateAclAuthorityName(t *testing.T) {
	keeper, _ := keepertest.AclKeeper(t)
	addr := sample.AccAddress()
	testCases := []struct {
		name     string
		input    types.AclAuthority
		newName  string
		expected string
	}{
		{
			name: "trim spaces around name",
			input: types.AclAuthority{
				Address: addr,
				Name:    "Old Name",
			},
			newName:  "  New Authority Name  ",
			expected: "New Authority Name",
		},
		{
			name: "empty string as name",
			input: types.AclAuthority{
				Address: addr,
				Name:    "Old Name",
			},
			newName:  "",
			expected: "",
		},
		{
			name: "name with no extra spaces",
			input: types.AclAuthority{
				Address: addr,
				Name:    "Old Name",
			},
			newName:  "CleanName",
			expected: "CleanName",
		},
		{
			name: "name with only whitespace",
			input: types.AclAuthority{
				Address: addr,
				Name:    "Old Name",
			},
			newName:  "     ",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aclAuthority := keeper.UpdateAclAuthorityName(tc.input, tc.newName)
			require.Equal(t, tc.expected, aclAuthority.Name)
		})
	}
}

func TestOverwriteAccessDefinitionsList(t *testing.T) {
	keeper, _ := keepertest.AclKeeper(t)
	addr := sample.AccAddress()
	testCases := []struct {
		name                    string
		inputAclAuthority       types.AclAuthority
		accessDefinitionListStr string
		expectedOutput          []*types.AccessDefinition
		expectedLen             int
		expErr                  bool
		expErrMsg               string
	}{
		{
			name: "empty access definition list",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionListStr: `[]`,
			expErr:                  true,
			expErrMsg:               "access definition list is empty",
		},
		{
			name: "invalid access definitions format",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionListStr: `[{"module":"module4","is_maker":true "is_checker":true}]`,
			expErr:                  true,
			expErrMsg:               "invalid AccessDefinitionList format",
		},
		{
			name: "empty module name",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionListStr: `[{"module":"","is_maker":true,"is_checker":true}]`,
			expErr:                  true,
			expErrMsg:               "invalid module name",
		},
		{
			name: "duplicated module",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionListStr: `[{"module":"module1","is_maker":true,"is_checker":true},{"module":"module1","is_maker":true,"is_checker":true}]`,
			expErr:                  true,
			expErrMsg:               "invalid module name",
		},
		{
			name: "at least one of is_maker or is_checker must be true",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionListStr: `[{"module":"module5","is_maker":false,"is_checker":false}]`,
			expErr:                  true,
			expErrMsg:               "at least one of is_maker or is_checker must be true",
		},
		{
			name: "all good",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionListStr: `[{"module":"module2","is_maker":true,"is_checker":true},{"module":"module3","is_maker":false,"is_checker":true}]`,
			expectedOutput: []*types.AccessDefinition{
				{Module: "module2", IsMaker: true, IsChecker: true},
				{Module: "module3", IsMaker: false, IsChecker: true},
			},
			expectedLen: 2,
			expErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aclAuthority, err := keeper.OverwriteAccessDefinitionList(tc.inputAclAuthority, tc.accessDefinitionListStr)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedLen, len(aclAuthority.AccessDefinitions))
				require.Equal(t, tc.expectedOutput, aclAuthority.AccessDefinitions)
			}
		})
	}
}

func TestUpdateAccessDefinitions(t *testing.T) {
	keeper, _ := keepertest.AclKeeper(t)
	addr := sample.AccAddress()
	testCases := []struct {
		name                       string
		inputAclAuthority          types.AclAuthority
		singleAccessDefinitionsStr string
		expectedOutput             []*types.AccessDefinition
		expErr                     bool
		expErrMsg                  string
	}{
		{
			name: "invalid access definitions format",
			inputAclAuthority: types.AclAuthority{
				Address:           addr,
				Name:              "Alice",
				AccessDefinitions: []*types.AccessDefinition{},
			},
			singleAccessDefinitionsStr: `{"module":"module1","is_maker":true "is_checker":false}`,
			expErr:                     true,
			expErrMsg:                  "invalid AccessDefinitionObject format",
		},
		{
			name: "update empty module",
			inputAclAuthority: types.AclAuthority{
				Address:           addr,
				Name:              "Alice",
				AccessDefinitions: []*types.AccessDefinition{},
			},
			singleAccessDefinitionsStr: `{"module":"","is_maker":true ,"is_checker":false}`,
			expErr:                     true,
			expErrMsg:                  "invalid module name",
		},
		{
			name: "at least one of is_maker or is_checker must be true",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: false, IsChecker: false},
				},
			},
			singleAccessDefinitionsStr: `{"module":"module1","is_maker":false ,"is_checker":false}`,
			expErr:                     true,
			expErrMsg:                  "at least one of is_maker or is_checker must be true",
		},
		{
			name: "fail when module does not exist in current ACL list",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: false, IsChecker: false},
				},
			},
			singleAccessDefinitionsStr: `{"module":"module2","is_maker":true,"is_checker":false}`,
			expErr:                     true,
			expErrMsg:                  "module not exist",
		},
		{
			name: "all good",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			singleAccessDefinitionsStr: `{"module":"module1","is_maker":true, "is_checker":false}`,
			expectedOutput: []*types.AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aclAuthority, err := keeper.UpdateAccessDefinitions(tc.inputAclAuthority, tc.singleAccessDefinitionsStr)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, aclAuthority.AccessDefinitions, tc.expectedOutput)
			}
		})
	}
}

func TestAddAccessDefinitions(t *testing.T) {
	keeper, _ := keepertest.AclKeeper(t)
	addr := sample.AccAddress()
	testCases := []struct {
		name                     string
		inputAclAuthority        types.AclAuthority
		accessDefinitionsListStr string
		expectedOutput           []*types.AccessDefinition
		expectedLen              int
		expErr                   bool
		expErrMsg                string
	}{
		{
			name: "empty access definition list",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionsListStr: `[]`,
			expErr:                   true,
			expErrMsg:                "access definition list is empty",
		},
		{
			name: "invalid access definitions format",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionsListStr: `[{"module":"module4","is_maker":true "is_checker":true}]`,
			expErr:                   true,
			expErrMsg:                "invalid AccessDefinitionList format",
		},
		{
			name: "empty module name",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionsListStr: `[{"module":"","is_maker":true,"is_checker":true}]`,
			expErr:                   true,
			expErrMsg:                "invalid module name",
		},
		{
			name: "duplicated module",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionsListStr: `[{"module":"module1","is_maker":true,"is_checker":true},{"module":"module1","is_maker":true,"is_checker":true}]`,
			expErr:                   true,
			expErrMsg:                "invalid module name",
		},
		{
			name: "at least one of is_maker or is_checker must be true",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			accessDefinitionsListStr: `[{"module":"module5","is_maker":false,"is_checker":false}]`,
			expErr:                   true,
			expErrMsg:                "at least one of is_maker or is_checker must be true",
		},
		{
			name: "add existing module",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			accessDefinitionsListStr: `[{"module":"module1","is_maker":true,"is_checker":true}]`,
			expErr:                   true,
			expErrMsg:                "module already exist",
		},
		{
			name: "all good",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			accessDefinitionsListStr: `[{"module":"module2","is_maker":true,"is_checker":true},{"module":"module3","is_maker":true,"is_checker":true}]`,
			expectedOutput: []*types.AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
				{Module: "module2", IsMaker: true, IsChecker: true},
				{Module: "module3", IsMaker: true, IsChecker: true},
			},
			expectedLen: 3,
			expErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aclAuthority, err := keeper.AddAccessDefinitions(tc.inputAclAuthority, tc.accessDefinitionsListStr)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedLen, len(aclAuthority.AccessDefinitions))
				require.Equal(t, tc.expectedOutput, aclAuthority.AccessDefinitions)
			}
		})
	}
}

func TestDeleteAccessDefinitions(t *testing.T) {
	keeper, _ := keepertest.AclKeeper(t)
	addr := sample.AccAddress()
	testCases := []struct {
		name              string
		inputAclAuthority types.AclAuthority
		moduleNames       []string
		expectedOutput    []*types.AccessDefinition
		expectedLen       int
		expErr            bool
		expErrMsg         string
	}{
		{
			name: "empty list",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			moduleNames: []string{},
			expErr:      true,
			expErrMsg:   "invalid module name",
		},
		{
			name: "empty access definition list",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
			},
			moduleNames: []string{"module1"},
			expErr:      true,
			expErrMsg:   "access definition list is empty",
		},

		{
			name: "module not found",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			moduleNames: []string{"module2"},
			expectedLen: 0,
			expErr:      true,
			expErrMsg:   "module not exist",
		},
		{
			name: "all good",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
					{Module: "module2", IsMaker: true, IsChecker: false},
					{Module: "module3", IsMaker: true, IsChecker: false},
					{Module: "module4", IsMaker: true, IsChecker: false},
				},
			},
			moduleNames: []string{"module1", "module2"},
			expectedOutput: []*types.AccessDefinition{
				{Module: "module3", IsMaker: true, IsChecker: false},
				{Module: "module4", IsMaker: true, IsChecker: false},
			},
			expectedLen: 2,
			expErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aclAuthority, err := keeper.DeleteAccessDefinitions(tc.inputAclAuthority, tc.moduleNames)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedLen, len(aclAuthority.AccessDefinitions))
				require.Equal(t, tc.expectedOutput, aclAuthority.AccessDefinitions)
			}
		})
	}
}

func TestClearAllAccessDefinitions(t *testing.T) {
	keeper, _ := keepertest.AclKeeper(t)
	addr := sample.AccAddress()
	testCases := []struct {
		name              string
		inputAclAuthority types.AclAuthority
		expectedLen       int
	}{
		{
			name: "clear single access definition",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			expectedLen: 0,
		},
		{
			name: "clear multiple access definitions",
			inputAclAuthority: types.AclAuthority{
				Address: addr,
				Name:    "Alice",
				AccessDefinitions: []*types.AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
					{Module: "module2", IsMaker: false, IsChecker: true},
				},
			},
			expectedLen: 0,
		},
		{
			name: "clear empty access definitions",
			inputAclAuthority: types.AclAuthority{
				Address:           addr,
				Name:              "Alice",
				AccessDefinitions: []*types.AccessDefinition{},
			},
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aclAuthority := keeper.ClearAllAccessDefinitions(tc.inputAclAuthority)
			require.Equal(t, tc.expectedLen, len(aclAuthority.AccessDefinitions))
		})
	}
}
