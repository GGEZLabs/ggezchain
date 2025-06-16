package e2e

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	tradetypes "github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testTrade() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	admin1, err := s.chainA.genesisAccounts[0].keyInfo.GetAddress() // receiver and signer
	s.Require().NoError(err)

	admin2, err := s.chainA.genesisAccounts[1].keyInfo.GetAddress()
	s.Require().NoError(err)

	fees := sdk.NewCoin(uggez1Denom, math.NewInt(1))

	// Create trade
	s.execCreateTrade(s.chainA, 0, admin1.String(), tradetypes.GetSampleTradeData(tradetypes.TradeTypeBuy), `{}`, `{}`, `{}`, admin1.String(), ggezHomePath, fees.String())

	s.Require().Eventually(
		func() bool {
			storedTrade, err := queryStoredTrade(chainEndpoint, "1")
			s.Require().NoError(err)

			storedTrades, err := queryAllStoredTrade(chainEndpoint)
			s.Require().NoError(err)

			storedTempTrade, err := queryStoredTempTrade(chainEndpoint, "1")
			s.Require().NoError(err)

			storedTempTrades, err := queryAllStoredTempTrade(chainEndpoint)
			s.Require().NoError(err)

			s.Require().Equal(storedTrade.StoredTrade.Status, tradetypes.StatusPending)
			s.Require().Equal(storedTempTrade.StoredTempTrade.TradeIndex, uint64(1))
			s.Require().Len(storedTempTrades.StoredTempTrade, 1)

			return len(storedTrades.StoredTrade) == 1
		},
		20*time.Second,
		5*time.Second,
	)

	// Process trade
	s.execProcessTrade(s.chainA, 0, "1", "confirm", admin2.String(), ggezHomePath, fees.String(), false)

	s.Require().Eventually(
		func() bool {
			storedTrade, err := queryStoredTrade(chainEndpoint, "1")
			s.Require().NoError(err)

			storedTempTrades, err := queryAllStoredTempTrade(chainEndpoint)
			s.Require().NoError(err)

			s.Require().Equal(storedTrade.StoredTrade.Status, tradetypes.StatusProcessed)

			ggzSupply, err := querySupplyOf(chainEndpoint, tradetypes.DefaultDenom)
			s.Require().NoError(err)

			s.Require().Equal(int64(1000000), ggzSupply.Amount.Int64())

			receiverAddressBalance, err := getSpecificBalance(chainEndpoint, admin1.String(), tradetypes.DefaultDenom)
			s.Require().NoError(err)

			s.Require().Equal(int64(1000000), receiverAddressBalance.Amount.Int64())

			return len(storedTempTrades.StoredTempTrade) == 0
		},
		20*time.Second,
		5*time.Second,
	)

	// Process already processed trade
	s.execProcessTrade(s.chainA, 0, "1", "confirm", admin2.String(), ggezHomePath, fees.String(), true)
}
