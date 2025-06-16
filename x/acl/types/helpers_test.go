package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAccessDefinitionList(t *testing.T) {
	tests := []struct {
		name                    string
		accessDefinitionListStr string
		expectedOutput          []*AccessDefinition
		expectedLen             int
		err                     error
	}{
		{
			name:                    "invalid access definitions format",
			accessDefinitionListStr: `[{"module":"module1","is_maker":true "is_checker":true}]`,
			err:                     ErrInvalidAccessDefinitionList,
		},
		{
			name:                    "empty access definition list",
			accessDefinitionListStr: `[]`,
			err:                     ErrEmptyAccessDefinitionList,
		},
		{
			name:                    "add empty module",
			accessDefinitionListStr: `[{"module":"","is_maker":true,"is_checker":true}]`,
			err:                     ErrInvalidModuleName,
		},
		{
			name:                    "add duplicated modules",
			accessDefinitionListStr: `[{"module":"module1","is_maker":true,"is_checker":true},{"module":"module1","is_maker":true,"is_checker":true}]`,
			err:                     ErrInvalidModuleName,
		},
		{
			name:                    "all good",
			accessDefinitionListStr: `[{"module":"module1","is_maker":false,"is_checker":true},{"module":"module2","is_maker":true,"is_checker":true}]`,
			expectedOutput: []*AccessDefinition{
				{Module: "module1", IsMaker: false, IsChecker: true},
				{Module: "module2", IsMaker: true, IsChecker: true},
			},
			expectedLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessDefinitions, err := ValidateAccessDefinitionList(tt.accessDefinitionListStr)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Len(t, accessDefinitions, tt.expectedLen)
			require.Equal(t, tt.expectedOutput, accessDefinitions)
		})
	}
}

func TestValidateSingleAccessDefinition(t *testing.T) {
	tests := []struct {
		name                string
		accessDefinitionStr string
		expectedOutput      *AccessDefinition
		err                 error
	}{
		{
			name:                "invalid access definitions format",
			accessDefinitionStr: `{"module":"module1","is_maker":true "is_checker":true}`,
			err:                 ErrInvalidAccessDefinitionObject,
		},
		{
			name:                "add empty module",
			accessDefinitionStr: `{"module":"","is_maker":true,"is_checker":true}`,
			err:                 ErrInvalidModuleName,
		},

		{
			name:                "all good",
			accessDefinitionStr: `{"module":"module1","is_maker":false,"is_checker":true}`,
			expectedOutput:      &AccessDefinition{Module: "module1", IsMaker: false, IsChecker: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessDefinition, err := ValidateSingleAccessDefinition(tt.accessDefinitionStr)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, accessDefinition)
		})
	}
}

func TestGetUpdatedAccessDefinitionList(t *testing.T) {
	tests := []struct {
		name           string
		currentList    []*AccessDefinition
		update         *AccessDefinition
		expectedOutput []*AccessDefinition
	}{
		{
			name: "update existing module (exact match)",
			currentList: []*AccessDefinition{
				{Module: "module1", IsMaker: false, IsChecker: false},
			},
			update: &AccessDefinition{Module: "module1", IsMaker: true, IsChecker: true},
			expectedOutput: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: true},
			},
		},
		{
			name: "update existing module (case insensitive match)",
			currentList: []*AccessDefinition{
				{Module: "module1", IsMaker: false, IsChecker: false},
			},
			update: &AccessDefinition{Module: "Module1", IsMaker: true, IsChecker: false},
			expectedOutput: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
			},
		},
		{
			name: "module not found (list unchanged)",
			currentList: []*AccessDefinition{
				{Module: "module1", IsMaker: false, IsChecker: false},
			},
			update: &AccessDefinition{Module: "module2", IsMaker: true, IsChecker: true},
			expectedOutput: []*AccessDefinition{
				{Module: "module1", IsMaker: false, IsChecker: false},
			},
		},
		{
			name:           "empty current list (no update)",
			currentList:    []*AccessDefinition{},
			update:         &AccessDefinition{Module: "module1", IsMaker: true, IsChecker: true},
			expectedOutput: []*AccessDefinition{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessDefinitionList := GetUpdatedAccessDefinitionList(tt.currentList, tt.update)
			require.Equal(t, tt.expectedOutput, accessDefinitionList)
		})
	}
}

func TestGetAuthorityModules(t *testing.T) {
	tests := []struct {
		name                 string
		accessDefinitionList []*AccessDefinition
		expectedOutput       []string
	}{
		{
			name:                 "empty list returns nil",
			accessDefinitionList: []*AccessDefinition{},
			expectedOutput:       nil,
		},
		{
			name: "single valid module",
			accessDefinitionList: []*AccessDefinition{
				{Module: "module1"},
			},
			expectedOutput: []string{"module1"},
		},
		{
			name: "multiple valid modules",
			accessDefinitionList: []*AccessDefinition{
				{Module: "module1"},
				{Module: "module2"},
				{Module: "module3"},
			},
			expectedOutput: []string{"module1", "module2", "module3"},
		},
		{
			name: "includes empty module (should be skipped)",
			accessDefinitionList: []*AccessDefinition{
				{Module: "module1"},
				{Module: ""},
				{Module: "module4"},
			},
			expectedOutput: []string{"module1", "module4"},
		},
		{
			name: "all modules empty (returns empty list)",
			accessDefinitionList: []*AccessDefinition{
				{Module: ""},
				{Module: ""},
			},
			expectedOutput: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessDefinitionList := GetAuthorityModules(tt.accessDefinitionList)
			require.Equal(t, tt.expectedOutput, accessDefinitionList)
		})
	}
}

func TestHasModule(t *testing.T) {
	tests := []struct {
		name           string
		currentModules []string
		module         string
		expectedOutput bool
	}{
		{
			name:           "module exists",
			currentModules: []string{"module1", "module2"},
			module:         "module1",
			expectedOutput: true,
		},
		{
			name:           "module not exists",
			currentModules: []string{"module1", "module2"},
			module:         "Module1",
			expectedOutput: false,
		},
		{
			name:           "module with additional spaces",
			currentModules: []string{"module1", "module2"},
			module:         "  module1  ",
			expectedOutput: false,
		},
		{
			name:           "empty module string",
			currentModules: []string{"module3", "module4"},
			module:         "",
			expectedOutput: false,
		},
		{
			name:           "nil module list",
			currentModules: nil,
			module:         "module4",
			expectedOutput: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasModule := hasModule(tt.currentModules, tt.module)
			require.Equal(t, tt.expectedOutput, hasModule)
		})
	}
}

func TestValidateUpdateAccessDefinition(t *testing.T) {
	tests := []struct {
		name          string
		aclAuthority  AclAuthority
		updatedModule string
		err           error
	}{
		{
			name: "module not exists",
			aclAuthority: AclAuthority{
				AccessDefinitions: []*AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			updatedModule: "module2",
			err:           ErrModuleDoesNotExist,
		},
		{
			name: "empty module",
			aclAuthority: AclAuthority{
				AccessDefinitions: []*AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			updatedModule: "",
			err:           ErrModuleDoesNotExist,
		},
		{
			name: "empty access definition list",
			aclAuthority: AclAuthority{
				AccessDefinitions: []*AccessDefinition{},
			},
			updatedModule: "module1",
			err:           ErrModuleDoesNotExist,
		},
		{
			name: "all good",
			aclAuthority: AclAuthority{
				AccessDefinitions: []*AccessDefinition{
					{Module: "module1", IsMaker: true, IsChecker: false},
				},
			},
			updatedModule: "module1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateAccessDefinition(tt.aclAuthority, tt.updatedModule)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateAddAccessDefinition(t *testing.T) {
	tests := []struct {
		name           string
		currentModules []string
		newModules     []string
		err            error
	}{
		{
			name:           "duplicate module",
			currentModules: []string{"module1", "module2"},
			newModules:     []string{"module1"},
			err:            ErrModuleExists,
		},
		{
			name:           "multiple duplicate modules",
			currentModules: []string{"module1", "module2", "module3"},
			newModules:     []string{"module1", "module3"},
			err:            ErrModuleExists,
		},
		{
			name:           "all good",
			currentModules: []string{},
			newModules:     []string{"module1", "module3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddAccessDefinition(tt.currentModules, tt.newModules)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateDeleteAccessDefinition(t *testing.T) {
	tests := []struct {
		name              string
		moduleNames       []string
		accessDefinitions []*AccessDefinition
		expectedOutput    []*AccessDefinition
		err               error
	}{
		{
			name:              "no module names",
			moduleNames:       []string{},
			accessDefinitions: []*AccessDefinition{},
			err:               ErrInvalidModuleName,
		},
		{
			name:              "access definitions list is empty",
			moduleNames:       []string{"module1"},
			accessDefinitions: []*AccessDefinition{},
			err:               ErrEmptyAccessDefinitionList,
		},
		{
			name:        "module to delete does not exist",
			moduleNames: []string{"module4"},
			accessDefinitions: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
				{Module: "module2", IsMaker: true, IsChecker: false},
				{Module: "module3", IsMaker: true, IsChecker: false},
			},
			err: ErrModuleDoesNotExist,
		},
		{
			name:        "remove one module",
			moduleNames: []string{"module1"},
			accessDefinitions: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
				{Module: "module2", IsMaker: true, IsChecker: false},
				{Module: "module3", IsMaker: true, IsChecker: false},
			},
			expectedOutput: []*AccessDefinition{
				{Module: "module2", IsMaker: true, IsChecker: false},
				{Module: "module3", IsMaker: true, IsChecker: false},
			},
		},
		{
			name:        "remove more than one module",
			moduleNames: []string{"module1", "module2"},
			accessDefinitions: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
				{Module: "module2", IsMaker: true, IsChecker: false},
				{Module: "module3", IsMaker: true, IsChecker: false},
			},
			expectedOutput: []*AccessDefinition{
				{Module: "module3", IsMaker: true, IsChecker: false},
			},
		},
		{
			name:        "remove all modules",
			moduleNames: []string{"module1", "module2", "module3"},
			accessDefinitions: []*AccessDefinition{
				{Module: "module1", IsMaker: true, IsChecker: false},
				{Module: "module2", IsMaker: true, IsChecker: false},
				{Module: "module3", IsMaker: true, IsChecker: false},
			},
			expectedOutput: []*AccessDefinition{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessDefinitions, err := ValidateDeleteAccessDefinition(tt.moduleNames, tt.accessDefinitions)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, accessDefinitions)
		})
	}
}

func TestValidateModuleOverlap(t *testing.T) {
	tests := []struct {
		name           string
		removedModules []string
		modules        []string
		err            error
	}{
		{
			name:           "single overlapping module",
			modules:        []string{"module1", "module2"},
			removedModules: []string{"module2"},
			err:            ErrUpdateAndRemoveModule,
		},
		{
			name:           "multiple overlapping module",
			modules:        []string{"module1", "module2"},
			removedModules: []string{"module1", "module2"},
			err:            ErrUpdateAndRemoveModule,
		},
		{
			name:           "all good",
			modules:        []string{"module1", "module2"},
			removedModules: []string{"module3", "module4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateModuleOverlap(tt.removedModules, tt.modules)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateAndExtractModuleNames(t *testing.T) {
	tests := []struct {
		name                string
		addAccessDefinition string
		expectedOutput      []string
		expectedLen         int
		err                 error
	}{
		{
			name:                "invalid access definitions format",
			addAccessDefinition: `[{"module":"module1","is_maker":true "is_checker":true}]`,
			err:                 ErrInvalidAccessDefinitionList,
		},
		{
			name:                "empty access definition list",
			addAccessDefinition: `[]`,
			err:                 ErrEmptyAccessDefinitionList,
		},
		{
			name:                "add empty module",
			addAccessDefinition: `[{"module":"","is_maker":true,"is_checker":true}]`,
			err:                 ErrInvalidModuleName,
		},
		{
			name:                "add duplicated modules",
			addAccessDefinition: `[{"module":"module1","is_maker":true,"is_checker":true},{"module":"module1","is_maker":true,"is_checker":true}]`,
			err:                 ErrInvalidModuleName,
		},
		{
			name:                "return one module",
			addAccessDefinition: `[{"module":"module1","is_maker":false,"is_checker":true}]`,
			expectedOutput:      []string{"module1"},
			expectedLen:         1,
		},
		{
			name:                "return more than one module",
			addAccessDefinition: `[{"module":"module1","is_maker":false,"is_checker":true},{"module":"module2","is_maker":true,"is_checker":true}]`,
			expectedOutput:      []string{"module1", "module2"},
			expectedLen:         2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules, err := validateAndExtractModuleNames(tt.addAccessDefinition)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, modules)
			require.Len(t, modules, tt.expectedLen)
		})
	}
}

func TestValidateConflictBetweenAccessDefinition(t *testing.T) {
	tests := []struct {
		name                   string
		updateAccessDefinition string
		addAccessDefinition    string
		removeList             []string
		err                    error
	}{
		{
			name:                   "invalid update access definitions format",
			updateAccessDefinition: `{"module":"module1","is_maker":true "is_checker":true}`,
			removeList:             []string{""},
			err:                    ErrInvalidAccessDefinitionObject,
		},
		{
			name:                   "update empty module",
			updateAccessDefinition: `{"module":"","is_maker":true,"is_checker":true}`,
			removeList:             []string{""},
			err:                    ErrInvalidModuleName,
		},
		{
			name:                   "update and remove same module",
			updateAccessDefinition: `{"module":"module1","is_maker":false,"is_checker":true}`,
			removeList:             []string{"module1"},
			err:                    ErrUpdateAndRemoveModule,
		},
		{
			name:                "invalid add access definitions format",
			addAccessDefinition: `[{"module":"module1","is_maker":true "is_checker":true}]`,
			removeList:          []string{""},
			err:                 ErrInvalidAccessDefinitionList,
		},
		{
			name:                "empty add access definition list",
			addAccessDefinition: `[]`,
			removeList:          []string{""},
			err:                 ErrEmptyAccessDefinitionList,
		},
		{
			name:                "add empty module",
			addAccessDefinition: `[{"module":"","is_maker":true,"is_checker":true}]`,
			removeList:          []string{""},
			err:                 ErrInvalidModuleName,
		},
		{
			name:                "add duplicated modules",
			addAccessDefinition: `[{"module":"module1","is_maker":true,"is_checker":true},{"module":"module1","is_maker":true,"is_checker":true}]`,
			removeList:          []string{""},
			err:                 ErrInvalidModuleName,
		},
		{
			name:                "add and remove same module",
			addAccessDefinition: `[{"module":"module1","is_maker":false,"is_checker":true}]`,
			removeList:          []string{"module1"},
			err:                 ErrUpdateAndRemoveModule,
		},
		{
			name:                   "add update and remove same module",
			addAccessDefinition:    `[{"module":"module1","is_maker":false,"is_checker":true}]`,
			updateAccessDefinition: `{"module":"module1","is_maker":false,"is_checker":true}`,
			removeList:             []string{"module1"},
			err:                    ErrUpdateAndRemoveModule,
		},
		{
			name:                   "all good",
			addAccessDefinition:    `[{"module":"module1","is_maker":false,"is_checker":true}]`,
			updateAccessDefinition: `{"module":"module1","is_maker":false,"is_checker":true}`,
			removeList:             []string{"module2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConflictBetweenAccessDefinition(tt.updateAccessDefinition, tt.addAccessDefinition, tt.removeList)
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateJSONFormat(t *testing.T) {
	tests := []struct {
		name      string
		jsonStr   string
		fieldName string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "invalid JSON format for field updateAccessDefinition",
			jsonStr:   `{"module":"module1","is_maker":true "is_checker":true}`,
			fieldName: "updateAccessDefinition",
			expErr:    true,
			expErrMsg: "invalid JSON format for field updateAccessDefinition",
		},
		{
			name:      "invalid JSON format for field addAccessDefinition",
			jsonStr:   `[{"module":"module1","is_maker":true "is_checker":true}]`,
			fieldName: "addAccessDefinition",
			expErr:    true,
			expErrMsg: "invalid JSON format for field addAccessDefinition",
		},
		{
			name:      "empty string",
			jsonStr:   "",
			fieldName: "updateAccessDefinition",
			expErr:    false,
		},
		{
			name:      "valid JSON format for field updateAccessDefinition",
			jsonStr:   `{"module":"module1","is_maker":true, "is_checker":true}`,
			fieldName: "updateAccessDefinition",
			expErr:    false,
		},
		{
			name:      "valid JSON format for field addAccessDefinition",
			jsonStr:   `[{"module":"module1","is_maker":true, "is_checker":true}]`,
			fieldName: "addAccessDefinition",
			expErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJSONFormat(tt.jsonStr, tt.fieldName)
			if tt.expErr {
				require.Contains(t, err.Error(), tt.expErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateDeletedModules(t *testing.T) {
	tests := []struct {
		name      string
		modules   []string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "empty module name",
			modules:   []string{""},
			expErr:    true,
			expErrMsg: "module name cannot be empty",
		},
		{
			name:      "white space module name",
			modules:   []string{"    "},
			expErr:    true,
			expErrMsg: "module name cannot be empty",
		},
		{
			name:      "duplicate module name",
			modules:   []string{"module1", "module1"},
			expErr:    true,
			expErrMsg: "module1 module(s) is duplicates",
		},
		{
			name:      "two duplicate module name",
			modules:   []string{"module1", "module1", "module2", "module2"},
			expErr:    true,
			expErrMsg: "module1, module2 module(s) is duplicates",
		},
		{
			name:    "all good",
			modules: []string{"module1", "module2"},
			expErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDeletedModules(tt.modules)
			if tt.expErr {
				require.Contains(t, err.Error(), tt.expErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
