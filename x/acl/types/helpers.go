package types

import (
	"encoding/json"
	"slices"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateAccessDefinitionList takes a JSON string of access definitions, validates it,
// and returns a structured slice of AccessDefinition or an error if invalid.
func ValidateAccessDefinitionList(accessDefinitionListStr string) ([]*AccessDefinition, error) {
	var accessDefinitionList []*AccessDefinition
	if err := json.Unmarshal([]byte(accessDefinitionListStr), &accessDefinitionList); err != nil {
		return nil, ErrInvalidAccessDefinitionList
	}

	if len(accessDefinitionList) == 0 {
		return nil, ErrEmptyAccessDefinitionList
	}

	seenModules := make(map[string]bool)
	var duplicateModules []string

	for _, accessDefinition := range accessDefinitionList {
		accessDefinition.Module = strings.ToLower(strings.TrimSpace(accessDefinition.Module))

		if accessDefinition.Module == "" {
			return nil, ErrInvalidModuleName.Wrapf("missing required parameter: module")
		}

		if seenModules[accessDefinition.Module] {
			duplicateModules = append(duplicateModules, accessDefinition.Module)
		}
		seenModules[accessDefinition.Module] = true
	}

	if len(duplicateModules) > 0 {
		return nil, ErrInvalidModuleName.Wrapf("duplicate module names found: %s", strings.Join(duplicateModules, ", "))
	}

	return accessDefinitionList, nil
}

// ValidateSingleAccessDefinition takes a JSON string for one AccessDefinition object, validates it,
// and returns the structured object or an error.
func ValidateSingleAccessDefinition(accessDefinitionStr string) (*AccessDefinition, error) {
	var accessDefinition AccessDefinition
	if err := json.Unmarshal([]byte(accessDefinitionStr), &accessDefinition); err != nil {
		return nil, ErrInvalidAccessDefinitionObject
	}

	accessDefinition.Module = strings.ToLower(strings.TrimSpace(accessDefinition.Module))

	if accessDefinition.Module == "" {
		return nil, ErrInvalidModuleName.Wrapf("missing required parameter: module")
	}
	return &accessDefinition, nil
}

// GetUpdatedAccessDefinitionList finds a module by name (case-insensitive)
// and updates its roles in-place within the current list.
func GetUpdatedAccessDefinitionList(currentList []*AccessDefinition, update *AccessDefinition) []*AccessDefinition {
	for _, m := range currentList {
		if strings.EqualFold(m.Module, update.Module) {
			m.IsMaker = update.IsMaker
			m.IsChecker = update.IsChecker
			break
		}
	}
	return currentList
}

// GetAuthorityModules retrieves the module names from a list of AccessDefinition.
func GetAuthorityModules(accessDefinitionList []*AccessDefinition) []string {
	if len(accessDefinitionList) == 0 {
		return nil
	}

	modules := make([]string, 0, len(accessDefinitionList))
	for _, accessDefinition := range accessDefinitionList {
		if accessDefinition.Module != "" {
			modules = append(modules, accessDefinition.Module)
		}
	}

	return modules
}

// ValidateUpdateAccessDefinition validate updated module
func ValidateUpdateAccessDefinition(aclAuthority AclAuthority, updatedModule string) error {
	modules := GetAuthorityModules(aclAuthority.AccessDefinitions)
	if !hasModule(modules, updatedModule) {
		return ErrModuleDoesNotExist.Wrapf("%s module not exist", updatedModule)
	}
	return nil
}

// ValidateAddAccessDefinition validates adding new modules, checking for existing duplicates.
func ValidateAddAccessDefinition(currentModules, newModules []string) error {
	existingModules := make(map[string]struct{}, len(currentModules))
	for _, module := range currentModules {
		existingModules[module] = struct{}{}
	}

	var duplicates []string
	for _, module := range newModules {
		if _, found := existingModules[module]; found {
			duplicates = append(duplicates, module)
		}
	}

	if len(duplicates) > 0 {
		return ErrModuleExists.Wrapf("%s module(s) already exist", strings.Join(duplicates, ", "))
	}

	return nil
}

// ValidateDeleteAccessDefinition validates module names for removal and returns the updated access list if valid.
func ValidateDeleteAccessDefinition(moduleNames []string, accessDefinitions []*AccessDefinition) ([]*AccessDefinition, error) {
	if len(moduleNames) == 0 {
		return nil, ErrInvalidModuleName.Wrapf("at least one module name must be provided")
	}

	if len(accessDefinitions) == 0 {
		return nil, ErrEmptyAccessDefinitionList
	}

	err := ValidateDeletedModules(moduleNames)
	if err != nil {
		return nil, err
	}

	modulesToRemove := make(map[string]struct{}, len(moduleNames))
	for _, name := range moduleNames {
		normalized := strings.ToLower(strings.TrimSpace(name))
		modulesToRemove[normalized] = struct{}{}
	}

	currentModules := make(map[string]struct{}, len(accessDefinitions))
	for _, module := range accessDefinitions {
		currentModules[module.Module] = struct{}{}
	}

	// Check if module not exist in current modules
	var missingModules []string
	for name := range modulesToRemove {
		if _, exists := currentModules[name]; !exists {
			missingModules = append(missingModules, name)
		}
	}

	if len(missingModules) > 0 {
		return nil, ErrModuleDoesNotExist.Wrapf("%s module(s) not found", strings.Join(missingModules, ", "))
	}

	updatedList := make([]*AccessDefinition, 0, len(accessDefinitions)-len(modulesToRemove))
	for _, module := range accessDefinitions {
		if _, shouldRemove := modulesToRemove[module.Module]; shouldRemove {
			continue
		}
		updatedList = append(updatedList, module)
	}

	return updatedList, nil
}

// ValidateConflictBetweenAccessDefinition validate update add and remove flags
func ValidateConflictBetweenAccessDefinition(updateAccessDefinition string, addAccessDefinition string, removeList []string) error {
	if updateAccessDefinition != "" && len(removeList) > 0 {
		updateAccessDefinition, err := ValidateSingleAccessDefinition(updateAccessDefinition)
		if err != nil {
			return err
		}

		if err = validateModuleOverlap([]string{updateAccessDefinition.Module}, removeList); err != nil {
			return err
		}
	}

	if addAccessDefinition != "" && len(removeList) > 0 {
		addedModules, err := validateAndExtractModuleNames(addAccessDefinition)
		if err != nil {
			return err
		}

		if err = validateModuleOverlap(addedModules, removeList); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDeletedModules checking for duplicates and empty modules
func ValidateDeletedModules(modules []string) error {
	seen := make(map[string]bool)
	var duplicates []string
	hasEmpty := false

	for _, module := range modules {
		normalized := strings.ToLower(strings.TrimSpace(module))
		if normalized == "" {
			hasEmpty = true
			continue
		}
		if seen[normalized] {
			duplicates = append(duplicates, module)
		} else {
			seen[normalized] = true
		}
	}

	switch {
	case hasEmpty:
		return ErrInvalidModuleName.Wrapf("module name cannot be empty")
	case len(duplicates) > 0:
		return ErrInvalidModuleName.Wrapf("%s module(s) is duplicates", strings.Join(duplicates, ", "))
	default:
		return nil
	}
}

// hasModule check if retrieved module exist in currentModules
func hasModule(currentModules []string, module string) bool {
	return slices.Contains(currentModules, module)
}

// validateModuleOverlap check if there is common modules between added or updated modules and removed modules
func validateModuleOverlap(modules, removedModules []string) error {
	removedSet := make(map[string]struct{}, len(removedModules))
	for _, module := range removedModules {
		removedSet[module] = struct{}{}
	}

	var overlaps []string
	// Check if any added/updated module exists in the removed set
	for _, module := range modules {
		if _, exists := removedSet[module]; exists {
			overlaps = append(overlaps, module)
		}
	}

	if len(overlaps) > 0 {
		return ErrUpdateAndRemoveModule.Wrapf("%s", strings.Join(overlaps, ", "))
	}
	return nil
}

// validateAndExtractModuleNames validates a AccessDefinition JSON string and extracts the list of module names.
func validateAndExtractModuleNames(addAccessDefinition string) ([]string, error) {
	accessDefinitionList, err := ValidateAccessDefinitionList(addAccessDefinition)
	if err != nil {
		return nil, err
	}
	modules := GetAuthorityModules(accessDefinitionList)
	return modules, nil
}

// validateJSONFormat validate JSON format
func validateJSONFormat(jsonStr string, fieldName string) error {
	if jsonStr != "" {
		if !json.Valid([]byte(jsonStr)) {
			return sdkerrors.ErrInvalidRequest.Wrapf("invalid JSON format for field %s", fieldName)
		}
	}
	return nil
}
