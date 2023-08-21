package routes

import (
	"mw-project/controller"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine{
	router := gin.Default()
	router.GET("/api/v1/wallet", controller.GetWallet)
	
	router.GET("/api/v1/wallet/transactions", controller.GetTransactions)
	router.POST("/api/v1/wallet", controller.EnableWallet)
	router.POST("/api/v1/init", controller.CreateWallet)
	router.POST("/api/v1/wallet/deposits", controller.DepostitWallet)
	router.POST("/api/v1/wallet/withdrawals", controller.WithdrawWallet)
	router.PATCH("/api/v1/wallet", controller.DisableWallet)
	return router
}