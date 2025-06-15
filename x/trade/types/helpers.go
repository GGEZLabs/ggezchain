package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateDate checks if the input date string is in the correct format and not after today's date.
func ValidateDate(blockTime time.Time, dateStr string) error {
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid date format: %s", err)
	}

	now := blockTime.Truncate(24 * time.Hour)
	parsedDate = parsedDate.Truncate(24 * time.Hour)

	if parsedDate.After(now) {
		return sdkerrors.ErrInvalidRequest.Wrapf("date cannot be in the future %s", parsedDate.Format(time.RFC3339))
	}

	return nil
}
