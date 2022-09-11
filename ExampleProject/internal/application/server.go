package application

import (
	"ExampleProject/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StartServer() {
	utils.InitializeLogger()

	router := gin.Default()
	errRouter := router.SetTrustedProxies([]string{"127.0.0.1"})
	if errRouter != nil {
		//log.Fatal("Log_SetTrustedProxies:", errRouter)
		utils.Logger.Fatal("Log_SetTrustedProxies:", zap.Error(errRouter))
	}

	utils.Logger.Info("connecting gin routes")
	UserRoute(router)
	utils.Logger.Info("connected")

	errRun := router.Run(":9090")
	if errRun != nil {
		utils.Logger.Fatal("log_router_run", zap.Error(errRun))
	}

	//token entegre
	//auth
	//register
	//login ve login ile erişilen endpoint
	//jwt
	//zap kullanılacak
}
