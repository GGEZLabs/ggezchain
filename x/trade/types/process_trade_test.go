package types

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestCheckerAndMakerNotTheSame(t *testing.T) {
	tests := []struct {
		name  string
		msg   MsgProcessTrade
		maker string
		err   string
	}{
		{
			name:  "maker and checker are not equals",
			msg:   MsgProcessTrade{Creator: "ggez1checkeraddress"},
			maker: "ggez1makeraddress",
		},
		{
			name:  "qual maker and checker",
			msg:   MsgProcessTrade{Creator: "ggez1sameaddress"},
			maker: "ggez1sameaddress",
			err:   ErrCheckerMustBeDifferent.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.validateCheckerIsNotMaker(tt.maker)

			if err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name   string
		status TradeStatus
		err    string
	}{
		{
			name:   "processed status",
			status: StatusProcessed,
			err:    ErrTradeStatusCompleted.Error(),
		},
		{
			name:   "rejected status",
			status: StatusRejected,
			err:    ErrTradeStatusRejected.Error(),
		},
		{
			name:   "canceled status",
			status: StatusCanceled,
			err:    ErrTradeStatusCanceled.Error(),
		},
		{
			name:   "pending status",
			status: StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &MsgProcessTrade{}

			err := msg.validateStatus(tt.status)

			if err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	checkerAdd := sample.AccAddress()
	makerAdd := sample.AccAddress()

	tests := []struct {
		name   string
		msg    MsgProcessTrade
		maker  string
		status TradeStatus
		err    error
	}{
		{
			name:   "process trade with invalid checker address",
			msg:    MsgProcessTrade{Creator: "invalid_address"},
			maker:  makerAdd,
			status: StatusPending,
			err:    ErrInvalidChecker,
		},
		{
			name:   "process trade with invalid checker address (empty)",
			msg:    MsgProcessTrade{Creator: ""},
			maker:  makerAdd,
			status: StatusPending,
			err:    ErrInvalidChecker,
		},
		{
			name:   "process trade with valid checker address",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  makerAdd,
			status: StatusPending,
			err:    nil,
		},
		{
			name:   "process trade with unspecified status",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  makerAdd,
			status: StatusNil,
			err:    ErrInvalidStatus,
		},
		{
			name:   "process trade with processed status",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  makerAdd,
			status: StatusProcessed,
			err:    ErrTradeStatusCompleted,
		},
		{
			name:   "process trade with Rejected status",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  makerAdd,
			status: StatusRejected,
			err:    ErrTradeStatusRejected,
		},
		{
			name:   "process trade with Canceled status",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  makerAdd,
			status: StatusCanceled,
			err:    ErrTradeStatusCanceled,
		},
		{
			name:   "process trade with Pending status",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  makerAdd,
			status: StatusPending,
		},
		{
			name:   "process trade with maker equal checker ",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  checkerAdd,
			status: StatusPending,
			err:    ErrCheckerMustBeDifferent,
		},
		{
			name:   "process trade with maker not equal checker ",
			msg:    MsgProcessTrade{Creator: checkerAdd},
			maker:  makerAdd,
			status: StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate(tt.status, tt.maker)
			if err != nil {
				require.Equal(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
