package api

import (
	"gr-blockchain-side/service/ledger"
	"gr-blockchain-side/service/portal"

	"github.com/gin-gonic/gin"
)

func SetUpSavingsAPI(router *gin.Engine) {
	router.POST("/ledger/createSavingsAccount", func(c *gin.Context) {
		ledger.CreateSavingsAccount(c)
	})

	router.POST("/ledger/settleSavingsAccount", func(c *gin.Context) {
		ledger.SettleSavingsAccount(c)
	})
}

func SetUpPortalAPI(router *gin.Engine) {
	router.POST("/portal/id", func(c *gin.Context) {
		portal.QueryTxnsBySavingsAccount(c)
	})

	router.POST("/portal/account", func(c *gin.Context) {
		portal.QueryTxnsByBankAccount(c)
	})
}
