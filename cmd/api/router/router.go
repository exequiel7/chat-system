package router

import (
	"chat-system/cmd/api/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type routerImpl struct {
	router *gin.Engine
	port   string
}

type Router interface {
	Setup()
	GetRouter() *gin.Engine
}

func NewRouter(port string) *routerImpl {
	return &routerImpl{
		router: gin.Default(),
		port:   port,
	}
}

// Router Setup
func (r routerImpl) configure(dbSession *gocql.Session) {
	r.router.Use(middlewares.CORSMiddleware())
	r.routes(dbSession)
	if err := r.router.Run(r.port); err != nil {
		err = fmt.Errorf("unable to start router error: %v(MISSING)", err)
		panic(err)
	}
}

func (r routerImpl) Setup(dbSession *gocql.Session) {
	r.configure(dbSession)
}

func (r routerImpl) GetRouter() *gin.Engine {
	return r.router
}
