package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jhonnyesquivel/quasar-op/docs"
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/server/handler/topsecret"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	topsecretService locate.TopSecretService
}

const VERSION_1 = "v1"

func New(host string, port uint, topSecretService locate.TopSecretService) Server {
	gin.SetMode("release")
	srv := Server{
		engine:           gin.New(),
		httpAddr:         fmt.Sprintf("%s:%d", host, port),
		topsecretService: topSecretService,
	}

	srv.registerRoutes()
	return srv
}

func (s *Server) Run() error {
	log.Println("Server running on", s.httpAddr)
	return s.engine.Run(s.httpAddr)
}

func (s *Server) registerRoutes() {
	// docs.SwaggerInfo.Title = "Operación Fuego de Quasar"
	// docs.SwaggerInfo.Description = "Operación Fuego de Quasar -  MELI challenge"
	// docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	//docs.SwaggerInfo.Schemes = []string{"https"}

	v1 := s.engine.Group("/api/v1")
	{
		v1.POST("/topsecret", topsecret.TopSecretPOSTHandler(s.topsecretService))
		v1.POST("/topsecret_split/:satellite", topsecret.TopSecretSplitPOSTHandler(s.topsecretService))
		v1.GET("/topsecret_split", topsecret.TopSecretSplitGETHandler(s.topsecretService))
	}
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
