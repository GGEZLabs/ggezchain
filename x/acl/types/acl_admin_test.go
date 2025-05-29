package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestValidateAddAdmin(t *testing.T) {
	alice := sample.AccAddress()
	testCases := []struct {
		name          string
		currentAdmins []AclAdmin
		newAdmins     []string
		expErr        bool
		expErrMsg     string
	}{
		{
			name: "address already exist",
			currentAdmins: []AclAdmin{
				{Address: alice},
				{Address: sample.AccAddress()},
				{Address: sample.AccAddress()},
			},
			newAdmins: []string{alice},
			expErr:    true,
			expErrMsg: "admin(s) already exist",
		},
		{
			name: "all good",
			currentAdmins: []AclAdmin{
				{Address: sample.AccAddress()},
				{Address: sample.AccAddress()},
			},
			newAdmins: []string{sample.AccAddress()},
			expErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAddAdmin(tc.currentAdmins, tc.newAdmins)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateDeleteAdmin(t *testing.T) {
	alice := sample.AccAddress()
	bob := sample.AccAddress()
	testCases := []struct {
		name          string
		currentAdmins []AclAdmin
		deletedAdmins []string
		expErr        bool
		expErrMsg     string
	}{
		{
			name: "delete last admin",
			currentAdmins: []AclAdmin{
				{Address: alice},
			},
			deletedAdmins: []string{alice},
			expErr:        true,
			expErrMsg:     "cannot delete all admins, at least one aclAdmin must remain",
		},
		{
			name: "delete all admins",
			currentAdmins: []AclAdmin{
				{Address: alice},
				{Address: bob},
			},
			deletedAdmins: []string{alice, bob},
			expErr:        true,
			expErrMsg:     "cannot delete all admins, at least one aclAdmin must remain",
		},
		{
			name: "address not exist",
			currentAdmins: []AclAdmin{
				{Address: sample.AccAddress()},
				{Address: sample.AccAddress()},
			},
			deletedAdmins: []string{sample.AccAddress(), sample.AccAddress()},
			expErr:        true,
			expErrMsg:     "admin(s) not exist",
		},
		{
			name: "all good",
			currentAdmins: []AclAdmin{
				{Address: alice},
				{Address: sample.AccAddress()},
			},
			deletedAdmins: []string{alice},
			expErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDeleteAdmin(tc.currentAdmins, tc.deletedAdmins)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestConvertStringsToAclAdmins(t *testing.T) {
	tests := []struct {
		name        string
		addresses   []string
		expectedLen int
	}{
		{
			name:        "empty addresses list",
			addresses:   []string{},
			expectedLen: 0,
		},
		{
			name:        "multiple addresses",
			addresses:   []string{sample.AccAddress(), sample.AccAddress()},
			expectedLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acladmins := ConvertStringsToAclAdmins(tt.addresses)
			require.Len(t, acladmins, tt.expectedLen)
		})
	}
}

func TestValidateAddresses(t *testing.T) {
	tests := []struct {
		name      string
		addresses []string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "invalid address",
			addresses: []string{sample.AccAddress(), "invalid_address"},
			expErr:    true,
			expErrMsg: "invalid address",
		},
		{
			name:      "empty addresses list",
			addresses: []string{},
			expErr:    false,
		},
		{
			name:      "valid addresses",
			addresses: []string{sample.AccAddress(), sample.AccAddress()},
			expErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAddresses(tt.addresses)
			if tt.expErr {
				require.Contains(t, err.Error(), tt.expErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHasDuplicateAddresses(t *testing.T) {
	duplicateAddress := sample.AccAddress()

	tests := []struct {
		name           string
		addresses      []string
		expectedOutput bool
	}{
		{
			name:           "duplicate addresses",
			addresses:      []string{duplicateAddress, duplicateAddress},
			expectedOutput: true,
		},
		{
			name:           "no duplicate addresses",
			addresses:      []string{sample.AccAddress(), sample.AccAddress()},
			expectedOutput: false,
		},
		{
			name:           "empty addresses list",
			addresses:      []string{},
			expectedOutput: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasDuplicate := hasDuplicateAddresses(tt.addresses)
			require.EqualValues(t, tt.expectedOutput, hasDuplicate)
		})
	}
}
