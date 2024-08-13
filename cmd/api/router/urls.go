package router

import (
	"chat-system/cmd/api/middlewares"
	"chat-system/cmd/api/router/factory"

	_ "chat-system/cmd/api/docs"

	"github.com/gocql/gocql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r routerImpl) routes(dbSession *gocql.Session) {
	factoryCtrl := factory.NewCtrlFactory(dbSession)

	r.router.GET("/ping", factoryCtrl.BuildPingController().Ping)

	users := r.router.Group("/users")
	{
		users.POST("/register", factoryCtrl.BuildUserController().RegisterUser)
		users.POST("/login", factoryCtrl.BuildUserController().VerifyUserPassword)
		users.GET("", factoryCtrl.BuildUserController().ListUsers)
	}

	messaging := r.router.Group("/messages")
	messaging.Use(middlewares.Auth())
	{
		messaging.POST("/send", factoryCtrl.BuildMessagingController().SendMessage)
		messaging.GET("/history/:senderID/:receiverID", factoryCtrl.BuildMessagingController().GetConversationHistory)
	}

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
