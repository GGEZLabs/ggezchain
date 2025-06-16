package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateAddAdmin validates whether any of the new admin addresses already exist
// in the current admin list
func ValidateAddAdmin(currentAdmins []AclAdmin, newAdmins []string) error {
	adminMap := make(map[string]struct{}, len(currentAdmins))
	for _, admin := range currentAdmins {
		adminMap[admin.Address] = struct{}{}
	}

	var duplicates []string

	for _, newAdmin := range newAdmins {
		if _, exists := adminMap[newAdmin]; exists {
			duplicates = append(duplicates, newAdmin)
		}
	}

	if len(duplicates) > 0 {
		return ErrAdminExist.Wrapf("%s", strings.Join(duplicates, ", "))
	}
	return nil
}

// ValidateDeleteAdmin ensures that deletion admins exist and that the admin list is not emptied.
func ValidateDeleteAdmin(currentAdmins []AclAdmin, deletedAdmins []string) error {
	adminMap := make(map[string]struct{}, len(currentAdmins))
	for _, admin := range currentAdmins {
		adminMap[admin.Address] = struct{}{}
	}

	var notExistingAdmins []string

	for _, deletedAdmin := range deletedAdmins {
		if _, exists := adminMap[deletedAdmin]; !exists {
			notExistingAdmins = append(notExistingAdmins, deletedAdmin)
		}
	}

	if len(notExistingAdmins) > 0 {
		return ErrAdminNotExist.Wrapf("%s", strings.Join(notExistingAdmins, ", "))
	}

	return nil
}

// ConvertStringsToAclAdmins converts a slice of strings (addresses) into a slice of AclAdmin structs
func ConvertStringsToAclAdmins(addresses []string) []AclAdmin {
	admins := make([]AclAdmin, len(addresses))
	for i, addr := range addresses {
		admins[i] = AclAdmin{Address: addr}
	}
	return admins
}

// validateAddresses validates whether any of the addresses not valid
func validateAddresses(addresses []string) error {
	for _, address := range addresses {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return fmt.Errorf("invalid address %s: %v", address, err)
		}
	}
	return nil
}

// hasDuplicateAddresses check if addresses list has a duplicate addresses
func hasDuplicateAddresses(addresses []string) bool {
	addressesMap := make(map[string]bool)

	for _, addr := range addresses {
		if addressesMap[addr] {
			return true
		}
		addressesMap[addr] = true
	}

	return false
}
