package e2e

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testBankTokenTransfer() {
	s.Run("send_tokens_between_accounts", func() {
		var (
			err           error
			valIdx        = 0
			c             = s.chainA
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)

		// define one sender and two recipient accounts
		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()
		charlie, _ := c.genesisAccounts[3].keyInfo.GetAddress()

		var beforeAliceUGGEZ1Balance,
			beforeBobUGGEZ1Balance,
			beforeCharlieUGGEZ1Balance,
			afterAliceUGGEZ1Balance,
			afterBobUGGEZ1Balance,
			afterCharlieUGGEZ1Balance sdk.Coin

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, alice.String(), uggez1Denom)
				s.Require().NoError(err)

				beforeBobUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, bob.String(), uggez1Denom)
				s.Require().NoError(err)

				beforeCharlieUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, charlie.String(), uggez1Denom)
				s.Require().NoError(err)

				return beforeAliceUGGEZ1Balance.IsValid() && beforeBobUGGEZ1Balance.IsValid() && beforeCharlieUGGEZ1Balance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		// alice sends tokens to bob
		s.execBankSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String(), standardFees.String(), false)

		// check that the transfer was successful
		s.Require().Eventually(
			func() bool {
				afterAliceUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, alice.String(), uggez1Denom)
				s.Require().NoError(err)

				afterBobUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, bob.String(), uggez1Denom)
				s.Require().NoError(err)

				// gasFeesBurnt := standardFees.Sub(sdk.NewCoin(uggez1Denom, math.NewInt(1000)))
				// alice's balance should be decremented by the token amount and the gas fees
				// if the difference between expected and actual balance is less than 500, consider it as a success
				// any small change in operation/code can result in the gasFee difference
				// we set the threshold to 500 to avoid false negatives
				expectedAfterAliceUGGEZ1Balance := beforeAliceUGGEZ1Balance.Sub(tokenAmount).Sub(standardFees)
				decremented := afterAliceUGGEZ1Balance.Sub(expectedAfterAliceUGGEZ1Balance).Amount.LT(math.NewInt(500))

				incremented := beforeBobUGGEZ1Balance.Add(tokenAmount).IsEqual(afterBobUGGEZ1Balance)

				return decremented && incremented
			},
			10*time.Second,
			5*time.Second,
		)

		// save the updated account balances of alice and bob
		beforeAliceUGGEZ1Balance, beforeBobUGGEZ1Balance = afterAliceUGGEZ1Balance, afterBobUGGEZ1Balance

		// alice sends tokens to bob and charlie, at once
		s.execBankMultiSend(s.chainA, valIdx, alice.String(), []string{bob.String(), charlie.String()}, tokenAmount.String(), standardFees.String(), false)

		s.Require().Eventually(
			func() bool {
				afterAliceUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, alice.String(), uggez1Denom)
				s.Require().NoError(err)

				afterBobUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, bob.String(), uggez1Denom)
				s.Require().NoError(err)

				afterCharlieUGGEZ1Balance, err = getSpecificBalance(chainEndpoint, charlie.String(), uggez1Denom)
				s.Require().NoError(err)

				// gasFeesBurnt := standardFees.Sub(sdk.NewCoin(uggez1Denom, math.NewInt(1016)))
				// alice's balance should be decremented by the token amount and the gas fees
				// if the difference between expected and actual balance is less than 500, consider it as a success
				// any small change in operation/code can result in the gasFee difference
				// we set the threshold to 500 to avoid false negatives
				expectedAfterAliceUGGEZ1Balance := beforeAliceUGGEZ1Balance.Sub(tokenAmount).Sub(tokenAmount).Sub(standardFees)
				decremented := afterAliceUGGEZ1Balance.Sub(expectedAfterAliceUGGEZ1Balance).Amount.LT(math.NewInt(500))

				incremented := beforeBobUGGEZ1Balance.Add(tokenAmount).IsEqual(afterBobUGGEZ1Balance) &&
					beforeCharlieUGGEZ1Balance.Add(tokenAmount).IsEqual(afterCharlieUGGEZ1Balance)

				return decremented && incremented
			},
			10*time.Second,
			5*time.Second,
		)
	})
}
