package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		AclAuthorityMap: []AclAuthority{}, AclAdminMap: []AclAdmin{}, SuperAdmin: nil}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	aclAuthorityIndexMap := make(map[string]struct{})

	for _, elem := range gs.AclAuthorityMap {
		index := fmt.Sprint(elem.Address)
		if _, ok := aclAuthorityIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for aclAuthority")
		}
		aclAuthorityIndexMap[index] = struct{}{}
	}
	aclAdminIndexMap := make(map[string]struct{})

	for _, elem := range gs.AclAdminMap {
		index := fmt.Sprint(elem.Address)
		if _, ok := aclAdminIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for aclAdmin")
		}
		aclAdminIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
