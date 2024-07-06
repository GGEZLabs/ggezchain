package types_test

import (
	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateProcess(t *testing.T) {
	checAdd := sample.AccAddress()
	makerAdd := sample.AccAddress()
	processTrade := types.MsgProcessTrade{}

	tests := []struct {
		name    string
		maker   string
		checker string
		status  string
		err     error
	}{
		{
			name:    "process trade with invalid checker address",
			checker: "xxxx",
			maker:   makerAdd,
			status:  "Pending",
			err:     types.ErrInvalidChecker,
		}, {
			name:    "process trade with invalid checker address (empty)",
			checker: "",
			maker:   makerAdd,
			status:  "Pending",
			err:     types.ErrInvalidChecker,
		}, {
			name:    "process trade with valid checker address",
			checker: checAdd,
			maker:   makerAdd,
			status:  "Pending",
			err:     nil,
		}, {
			name:    "process trade with invalid status",
			checker: checAdd,
			maker:   makerAdd,
			status:  "xxxx",
			err:     types.ErrInvalidStatus,
		}, {
			name:    "process trade with empty status",
			checker: checAdd,
			maker:   makerAdd,
			status:  "",
			err:     types.ErrInvalidStatus,
		}, {
			name:    "process trade with Completed status",
			checker: checAdd,
			maker:   makerAdd,
			status:  "Completed",
			err:     types.ErrTradeStatusCompleted,
		}, {
			name:    "process trade with Rejected status",
			checker: checAdd,
			maker:   makerAdd,
			status:  "Rejected",
			err:     types.ErrTradeStatusRejected,
		}, {
			name:    "process trade with Canceled status",
			checker: checAdd,
			maker:   makerAdd,
			status:  "Canceled",
			err:     types.ErrTradeStatusCanceled,
		}, {
			name:    "process trade with Pending status",
			checker: checAdd,
			maker:   makerAdd,
			status:  "Pending",
			err:     nil,
		}, {
			name:    "process trade with maker equal checker ",
			checker: checAdd,
			maker:   checAdd,
			status:  "Pending",
			err:     types.ErrCheckerMustBeDifferent,
		}, {
			name:    "process trade with maker not equal checker ",
			checker: checAdd,
			maker:   makerAdd,
			status:  "Pending",
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processTrade.ValidateProcess(tt.status, tt.maker, tt.checker)
			if err != nil {
				require.Equal(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}

}
