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
		AclAuthorities: []AclAuthority{},
		AclAdmins:      []AclAdmin{},
		SuperAdmin:     nil,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if gs.SuperAdmin == nil {
		if len(gs.AclAdmins) > 0 || len(gs.AclAuthorities) > 0 {
			return fmt.Errorf("cannot initialize admins or authorities without a super admin: super admin must be set")
		}
		return nil
	}

	if len(gs.AclAdmins) == 0 && len(gs.AclAuthorities) > 0 {
		return fmt.Errorf("cannot initialize authorities without admin: admin must be set")
	}

	err := gs.ValidateSuperAdmin()
	if err != nil {
		return err
	}

	err = gs.ValidateAclAdmin()
	if err != nil {
		return err
	}

	err = gs.ValidateAclAuthority()
	if err != nil {
		return err
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}

func (gs GenesisState) ValidateAclAuthority() error {
	aclAuthorityIndexMap := make(map[string]struct{})

	for _, authority := range gs.AclAuthorities {
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
		}
	}
	return nil
}

func (gs GenesisState) ValidateAclAdmin() error {
	aclAdminIndexMap := make(map[string]struct{})

	for _, elem := range gs.AclAdmins {
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

func (gs GenesisState) ValidateSuperAdmin() error {
	if gs.SuperAdmin == nil {
		return nil
	}

	if _, err := sdk.AccAddressFromBech32(gs.SuperAdmin.Admin); err != nil {
		return fmt.Errorf("invalid super admin address %s: %w", gs.SuperAdmin.Admin, err)
	}

	return nil
}
