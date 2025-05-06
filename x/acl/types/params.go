package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyAdmin            = []byte("Admin")
	DefaultAdmin string = "ggez12mzusset7qzp2ndvm3ya9ss7vzucvfcj4f6npy"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	admin string,
) Params {
	return Params{
		Admin: admin,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdmin,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdmin, &p.Admin, validateAdmin),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateAdmin(p.Admin); err != nil {
		return err
	}

	return nil
}

// validateAdmin validates the Admin param
func validateAdmin(v interface{}) error {
	admin, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if admin == "" {
		return fmt.Errorf("admin address cannot be empty")
	}

	_, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return fmt.Errorf("invalid admin address: %w", err)
	}
	
	return nil
}
