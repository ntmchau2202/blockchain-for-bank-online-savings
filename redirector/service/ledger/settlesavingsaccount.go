package ledger

import (
	"gr-blockchain-side/internal/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SettleSavingsAccount(c *gin.Context) {
	// parsing
	var msg client.SettleSavingsAccountMessage
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
		return
	}

	transaction, err := bcClient.SaveSavingsAccountSettlement(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, client.CreateErrorMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, client.CreateSuccessMessage(transaction))
}
