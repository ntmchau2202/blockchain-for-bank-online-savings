package ledger

import (
	"gr-blockchain-side/api/endpoint"
	"gr-blockchain-side/internal/client"
	"gr-blockchain-side/internal/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SettleSavingsAccount(c *gin.Context) {
	// parsing
	var msg message.Request
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	if !msg.ValidateCommand(message.SearchTxnsBySavingsAccountCmd) {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage("command mismatch"))
		return
	}

	// customerID := msg.Details["customer_id"]
	// customerPhone := msg.Details["customer_phone"]
	// savingsAcountID := msg.Details["savingsaccount_id"]
	actualInterestAmount := msg.Details["actual_interest_amount"]
	settleTime := msg.Details["settle_time"]
	// bankName := msg.Details["bank_name"]

	if err := endpoint.ValidateAmount(actualInterestAmount); err != nil {
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

	transaction, err := bcClient.SaveSavingsAccountSettlement(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	// TODO: fill in the proper arguments here
	c.JSON(http.StatusOK, client.CreateSuccessMessage(transaction))
}
