package http

import "github.com/gin-gonic/gin"

type ApiServer struct {
	server *gin.Engine
	router *gin.RouterGroup
}
type Controller interface {
	SetupRoutes(r gin.IRouter)
}

func NewApiServer() *ApiServer {
	r := gin.Default()
	g := r.Group("api")

	return &ApiServer{
		server: r,
		router: g,
	}
}

func (s *ApiServer) SetupControllers(controllers []Controller) {
	for _, controller := range controllers {
		controller.SetupRoutes(s.router)
	}
}

func (s *ApiServer) Start() error {
	err := s.server.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
