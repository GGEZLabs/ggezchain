package types

import (
	"fmt"
	"strings"
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

// FormatPrice convert a float to a decimal string
func FormatPrice(price float64) string {
	str := fmt.Sprintf("%.12f", price)
	return strings.TrimRight(strings.TrimRight(str, "0"), ".")
}
