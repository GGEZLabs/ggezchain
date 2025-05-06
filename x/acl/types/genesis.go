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
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in aclAuthority
	aclAuthorityIndexMap := make(map[string]struct{})

	for _, elem := range gs.AclAuthorityList {

		// Validate address format
		if _, err := sdk.AccAddressFromBech32(elem.Address); err != nil {
			return fmt.Errorf("invalid address for aclAuthority: %s, error: %w", elem.Address, err)
		}

		index := string(AclAuthorityKey(elem.Address))
		if _, ok := aclAuthorityIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for aclAuthority")
		}
		aclAuthorityIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
