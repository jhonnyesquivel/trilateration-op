package bootstrap

import (
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/server"
)

const (
	host = "localhost"
	port = 8080
)

func Run() error {
	topsecretService := locate.NewTopSecretService()
	srv := server.New(host, port, topsecretService)
	return srv.Run()
}
