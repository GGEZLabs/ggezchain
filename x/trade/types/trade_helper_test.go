package types_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateReceiverAndCreatorAddress(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("validAddress")).String()
	invalidAddress := "invalidAddress123"

	tests := []struct {
		name     string
		receiver string
		creator  string
		err      string
	}{
		{
			name:     "valid receiver and creator address",
			receiver: validAddress,
			creator:  validAddress,
		},
		{
			name:     "invalid receiver address",
			receiver: invalidAddress,
			err:      types.ErrInvalidReceiverAddress.Error(),
		},
		{
			name:     "invalid creator address",
			receiver: validAddress,
			creator:  invalidAddress,
			err:      types.ErrInvalidCreator.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.MsgCreateTrade{
				ReceiverAddress: tt.receiver,
				Creator:         tt.creator,
			}

			err := msg.ValidateReceiverAndCreatorAddress()

			if err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCheckerAddress(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("validAddress")).String()
	invalidAddress := "invalidAddress123"

	tests := []struct {
		name    string
		checker string
		err     string
	}{
		{
			name:    "valid checker address",
			checker: validAddress,
		},
		{
			name:    "invalid checker address",
			checker: invalidAddress,
			err:     types.ErrInvalidChecker.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.MsgProcessTrade{}

			err := msg.ValidateCheckerAddress(tt.checker)

			if err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetPrepareCoin(t *testing.T) {
	storedTradeOne := types.StoredTrade{
		TradeIndex: 1,
		TradeType:  types.Buy,
		Coin:       types.DefaultCoinDenom,
		Price:      "10.59",
		Quantity:   "5026505",
	}
	storedTradeTwo := types.StoredTrade{
		TradeIndex: 2,
		TradeType:  types.Buy,
		Coin:       types.DefaultCoinDenom,
		Price:      "10.59",
		Quantity:   "invalid",
	}

	storedTradeThree := types.StoredTrade{
		TradeIndex: 3,
		TradeType:  types.Buy,
		Coin:       types.DefaultCoinDenom,
		Price:      "10.59",
		Quantity:   "92233720368547758077",
	}

	tests := []struct {
		name        string
		storedTrade types.StoredTrade
		err         string
	}{
		{
			name:        "valid quantity address",
			storedTrade: storedTradeOne,
		},
		{
			name:        "invalid quantity address",
			storedTrade: storedTradeTwo,
			err:         types.ErrInvalidTradeQuantity.Error(),
		},
		{
			name:        "too large quantity address",
			storedTrade: storedTradeThree,
			err:         types.ErrInvalidTradeQuantity.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.MsgProcessTrade{}

			_, err := msg.GetPrepareCoin(tt.storedTrade)

			if err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckerAndMakerNotTheSame(t *testing.T) {
	tests := []struct {
		name    string
		maker   string
		checker string
		err     string
	}{
		{
			name:    "equal maker and checker",
			maker:   "ggez1makeraddress",
			checker: "ggez1checkeraddress",
		},
		{
			name:    "maker and checker are not equals",
			maker:   "ggez1sameaddress",
			checker: "ggez1sameaddress",
			err:     types.ErrCheckerMustBeDifferent.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.MsgProcessTrade{}

			err := msg.CheckerAndMakerNotTheSame(tt.maker, tt.checker)

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
		status string
		err    string
	}{
		{
			name:   "completed status",
			status: types.Completed,
			err:    types.ErrTradeStatusCompleted.Error(),
		},
		{
			name:   "rejected status",
			status: types.Rejected,
			err:    types.ErrTradeStatusRejected.Error(),
		},
		{
			name:   "canceled status",
			status: types.Canceled,
			err:    types.ErrTradeStatusCanceled.Error(),
		},
		{
			name:   "pending status",
			status: types.Pending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.MsgProcessTrade{}

			err := msg.ValidateStatus(tt.status)

			if err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateProcess(t *testing.T) {
	checkerAdd := sample.AccAddress()
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
			checker: checkerAdd,
			maker:   makerAdd,
			status:  "Pending",
			err:     nil,
		}, {
			name:    "process trade with invalid status",
			checker: checkerAdd,
			maker:   makerAdd,
			status:  "xxxx",
			err:     types.ErrInvalidStatus,
		}, {
			name:    "process trade with empty status",
			checker: checkerAdd,
			maker:   makerAdd,
			status:  "",
			err:     types.ErrInvalidStatus,
		}, {
			name:    "process trade with Completed status",
			checker: checkerAdd,
			maker:   makerAdd,
			status:  "Completed",
			err:     types.ErrTradeStatusCompleted,
		}, {
			name:    "process trade with Rejected status",
			checker: checkerAdd,
			maker:   makerAdd,
			status:  "Rejected",
			err:     types.ErrTradeStatusRejected,
		}, {
			name:    "process trade with Canceled status",
			checker: checkerAdd,
			maker:   makerAdd,
			status:  "Canceled",
			err:     types.ErrTradeStatusCanceled,
		}, {
			name:    "process trade with Pending status",
			checker: checkerAdd,
			maker:   makerAdd,
			status:  "Pending",
			err:     nil,
		}, {
			name:    "process trade with maker equal checker ",
			checker: checkerAdd,
			maker:   checkerAdd,
			status:  "Pending",
			err:     types.ErrCheckerMustBeDifferent,
		}, {
			name:    "process trade with maker not equal checker ",
			checker: checkerAdd,
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
