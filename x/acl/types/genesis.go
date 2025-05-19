package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AclAuthorityList: []AclAuthority{},
		AclAdminList:     []AclAdmin{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateAclAuthority()
	if err != nil {
		return err
	}

	err = gs.ValidateAclAdmin()
	if err != nil {
		return err
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}

func (gs GenesisState) ValidateAclAuthority() error {
	aclAuthorityIndexMap := make(map[string]struct{})

	for _, authority := range gs.AclAuthorityList {
		// Validate address format
		if _, err := sdk.AccAddressFromBech32(authority.Address); err != nil {
			return fmt.Errorf("invalid address for aclAuthority: %s, error: %w", authority.Address, err)
		}

		// Validate name
		if authority.Name == "" {
			return fmt.Errorf("empty name not allowed, authority address: %s", authority.Address)
		}
		
		// Check for duplicated index in aclAuthority
		index := string(AclAuthorityKey(authority.Address))
		if _, ok := aclAuthorityIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for aclAuthority")
		}
		aclAuthorityIndexMap[index] = struct{}{}

		seenModules := make(map[string]struct{})

		// Validate access definitions
		for _, access := range authority.AccessDefinitions {
			// Check for duplicate modules
			if _, exists := seenModules[access.Module]; exists {
				return fmt.Errorf("duplicate module '%s' found in access definitions for address '%s'", access.Module, authority.Address)
			}
			seenModules[access.Module] = struct{}{}

			// todo
			// Must have at least one role
			// if !access.IsMaker && !access.IsChecker {
			// 	return fmt.Errorf("access definition for module '%s' must be either maker or checker (address: %s)", access.Module, authority.Address)
			// }
		}
	}
	return nil
}

func (gs GenesisState) ValidateAclAdmin() error {
	aclAdminIndexMap := make(map[string]struct{})

	for _, elem := range gs.AclAdminList {
		// Validate address format
		if _, err := sdk.AccAddressFromBech32(elem.Address); err != nil {
			return fmt.Errorf("invalid address for aclAdmin: %s, error: %w", elem.Address, err)
		}

		// Check for duplicated index in aclAdmin
		index := string(AclAdminKey(elem.Address))
		if _, ok := aclAdminIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for aclAdmin")
		}
		aclAdminIndexMap[index] = struct{}{}
	}
	return nil
}
