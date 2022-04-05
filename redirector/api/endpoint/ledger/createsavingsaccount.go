package ledger

import (
	"gr-blockchain-side/api/endpoint"
	"gr-blockchain-side/internal/blockchain/client"
	"gr-blockchain-side/internal/message"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSavingsAccount(c *gin.Context) {
	var msg message.Request
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	if !msg.ValidateCommand(message.SearchTxnsBySavingsAccountCmd) {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage("command mismatch"))
		return
	}

	savingsPeriod := msg.Details["savings_period"]
	interestRate := msg.Details["interest_rate"]
	estimatedInterestAmount := msg.Details["estimated_interest_amount"]
	settleTime := msg.Details["open_time"]
	currency := msg.Details["currency"]

	if err := endpoint.ValidateAmount(estimatedInterestAmount); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	if err := endpoint.ValidateCurrency(currency); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	if err := endpoint.ValidateInterestRate(interestRate); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	if err := endpoint.ValidateSavingsPeriod(savingsPeriod); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	if err := endpoint.ValidateTime(settleTime); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	bcClient, err := client.NewBlockchainClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	transaction, err := bcClient.SaveSavingsAccountCreation(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, client.CreateErrorMessage(err.Error()))
		log.Panic(err)
		return
	}

	c.JSON(http.StatusOK, client.CreateSuccessMessage(transaction))
}
