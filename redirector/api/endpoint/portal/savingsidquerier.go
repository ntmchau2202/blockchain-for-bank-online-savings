package portal

import (
	"gr-blockchain-side/api/endpoint"
	"gr-blockchain-side/internal/blockchain/client"
	"gr-blockchain-side/internal/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryTxnsBySavingsAccount(c *gin.Context) {
	var msg message.Request
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	if !msg.ValidateCommand(message.SearchTxnsBySavingsAccountCmd) {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage("command mismatch"))
		return
	}

	id := msg.Details["query_id"]
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, endpoint.CreateErrorMessage("missing query id"))
		return
	}

	bcClient, err := client.NewBlockchainClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	// TODO: reformulate the struct of message.Details
	transaction, err := bcClient.QueryTxnsBySavingsAccount(msg.Details["query_id"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, endpoint.CreateErrorMessage(err.Error()))
		return
	}

	// TODO: fill in this line the appropiate values

	// c.JSON(http.StatusOK, endpoint.CreateSuccessMessage("transaction", "", ""))
}
