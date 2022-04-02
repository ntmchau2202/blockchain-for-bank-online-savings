package portal

import (
	"gr-blockchain-side/internal/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryTxnsBySavingsAccount(c *gin.Context) {
	var msg client.QueryMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, client.CreateErrorMessage(err.Error()))
	}

	if err := msg.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, client.CreateErrorMessage(err.Error()))
		return
	}

	bcClient, err := client.NewBlockchainClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, client.CreateErrorMessage(err.Error()))
	}

	transaction, err := bcClient.QueryTxnsBySavingsAccount(msg.Details["query_id"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, client.CreateErrorMessage(err.Error()))
	}

	c.JSON(http.StatusOK, client.CreateSuccessMessage(transaction))
}
