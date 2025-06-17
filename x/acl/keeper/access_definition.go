package keeper

import (
	"strings"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
)

// UpdateAclAuthorityName update authority name
func (k Keeper) UpdateAclAuthorityName(aclAuthority types.AclAuthority, name string) types.AclAuthority {
	aclAuthority.Name = strings.TrimSpace(name)

	return aclAuthority
}

// OverwriteAccessDefinitionList completely replace the access definition list
func (k Keeper) OverwriteAccessDefinitionList(aclAuthority types.AclAuthority, accessDefinitionsListStr string) (types.AclAuthority, error) {
	overwriteAccessDefinitions, err := types.ValidateAccessDefinitionList(accessDefinitionsListStr)
	if err != nil {
		return types.AclAuthority{}, err
	}

	aclAuthority.AccessDefinitions = overwriteAccessDefinitions

	return aclAuthority, nil
}

// UpdateAccessDefinition update specific module permission
func (k Keeper) UpdateAccessDefinition(aclAuthority types.AclAuthority, singleAccessDefinitionsStr string) (types.AclAuthority, error) {
	updatedAccessDefinitions, err := types.ValidateSingleAccessDefinition(singleAccessDefinitionsStr)
	if err != nil {
		return types.AclAuthority{}, err
	}

	err = types.ValidateUpdateAccessDefinition(aclAuthority, updatedAccessDefinitions.Module)
	if err != nil {
		return types.AclAuthority{}, err
	}

	updatedAccessDefinitionsList := types.GetUpdatedAccessDefinitionList(aclAuthority.AccessDefinitions, updatedAccessDefinitions)
	aclAuthority.AccessDefinitions = updatedAccessDefinitionsList

	return aclAuthority, nil
}

// AddAccessDefinitions add one or more access definition
func (k Keeper) AddAccessDefinitions(aclAuthority types.AclAuthority, accessDefinitionsListStr string) (types.AclAuthority, error) {
	accessDefinitionsList, err := types.ValidateAccessDefinitionList(accessDefinitionsListStr)
	if err != nil {
		return types.AclAuthority{}, err
	}

	newModules := types.GetAuthorityModules(accessDefinitionsList)
	currentModules := types.GetAuthorityModules(aclAuthority.AccessDefinitions)

	err = types.ValidateAddAccessDefinition(currentModules, newModules)
	if err != nil {
		return types.AclAuthority{}, err
	}

	aclAuthority.AccessDefinitions = append(aclAuthority.AccessDefinitions, accessDefinitionsList...)

	return aclAuthority, nil
}

// DeleteAccessDefinitions remove one or more access definition
func (k Keeper) DeleteAccessDefinitions(aclAuthority types.AclAuthority, moduleNames []string) (types.AclAuthority, error) {
	newAccessDefinitionsList, err := types.ValidateDeleteAccessDefinition(moduleNames, aclAuthority.AccessDefinitions)
	if err != nil {
		return types.AclAuthority{}, err
	}

	aclAuthority.AccessDefinitions = newAccessDefinitionsList

	return aclAuthority, nil
}

// ClearAllAccessDefinitions clear all access definition
func (k Keeper) ClearAllAccessDefinitions(aclAuthority types.AclAuthority) types.AclAuthority {
	aclAuthority.AccessDefinitions = make([]*types.AccessDefinition, 0)

	return aclAuthority
}
