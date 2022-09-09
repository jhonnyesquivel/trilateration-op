package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/server/handler/topsecret"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	topsecretService locate.TopSecretService
}

const VERSION_1 = "v1"

func New(host string, port uint, topSecretService locate.TopSecretService) Server {
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
	s.engine.POST(fmt.Sprintf("%v/topsecret", VERSION_1), topsecret.TopSecretGETHandler(s.topsecretService))
	s.engine.POST(fmt.Sprintf("%v/topsecret_split/:satellite", VERSION_1), topsecret.TopSecretSplitPOSTHandler(s.topsecretService))
	s.engine.GET(fmt.Sprintf("%v/topsecret_split", VERSION_1), topsecret.TopSecretSplitGETHandler(s.topsecretService))
}
