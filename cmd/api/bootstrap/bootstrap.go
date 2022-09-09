package bootstrap

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/server"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/storage/pgsql"
)

type Configuration struct {
	Host  string
	Port  uint
	dbUrl string
}

var Config Configuration

func init() {

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 32)
	if err != nil {
		panic("PORT env var has an unvalid value")
	}
	Config.Host = os.Getenv("HOST")
	Config.dbUrl = os.Getenv("DATABASE_URL")
	Config.Port = uint(port)
}

func Run() error {

	creds, err := pg.ParseURL(Config.dbUrl)
	if err != nil {
		panic("db not ready")
	}

	db := pg.Connect(creds)
	db.AddQueryHook(dbLogger{})

	satelliteRepository := pgsql.NewSatelliteRepository(db)
	topsecretService := locate.NewTopSecretService(satelliteRepository)
	srv := server.New(Config.Host, Config.Port, topsecretService)
	return srv.Run()
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, e *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	fmt.Println(string(fq))
	return nil
}
