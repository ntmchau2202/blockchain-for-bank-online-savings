package ledger

import (
	"gr-blockchain-side/internal/client"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSavingsAccount(c *gin.Context) {
	// parsing
	var msg client.AddSavingsAccountMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, client.CreateErrorMessage(err.Error()))
		return
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
	transaction, err := bcClient.SaveSavingsAccountCreation(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, client.CreateErrorMessage(err.Error()))
		log.Panic(err)
		return
	}

	c.JSON(http.StatusOK, client.CreateSuccessMessage(transaction))
}
