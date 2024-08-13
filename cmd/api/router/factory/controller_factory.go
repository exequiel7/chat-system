package factory

import (
	"chat-system/cmd/api/controllers"
	"chat-system/cmd/api/repositories"
	"chat-system/cmd/api/services"

	"github.com/gocql/gocql"
)

type ControllerBuilder interface {
	BuildPingController() controllers.PingController
	BuildUserController() controllers.UserController
	BuildMessagingController() controllers.MessagingController
}

type controllerBuildImpl struct {
	dbSession *gocql.Session
}

func NewCtrlFactory(dbSession *gocql.Session) ControllerBuilder {
	return &controllerBuildImpl{dbSession: dbSession}
}

func (ctrlFactory *controllerBuildImpl) BuildPingController() controllers.PingController {
	return controllers.NewPingController()
}

func (ctrlFactory *controllerBuildImpl) BuildUserController() controllers.UserController {
	dbRepository := repositories.NewDBRepository(ctrlFactory.dbSession)
	userService := services.NewUserService(dbRepository)
	return controllers.NewUserController(userService)
}

func (ctrlFactory *controllerBuildImpl) BuildMessagingController() controllers.MessagingController {
	dbRepository := repositories.NewDBRepository(ctrlFactory.dbSession)
	msgService := services.NewMessagingService(dbRepository)
	return controllers.NewMessagingController(msgService)
}
