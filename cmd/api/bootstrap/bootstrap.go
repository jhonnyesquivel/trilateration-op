package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/server"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/storage/pgsql"
)

const (
	host   = "localhost"
	port   = 8080
	dbUser = "postgres"
	dbPass = "postgrespw"
	dbHost = "host.docker.internal"
	dbPort = "49153"
	dbName = "public"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, e *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	fmt.Println(string(fq))
	return nil
}

func Run() error {

	// mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	// db := pg.Connect(&pg.Options{
	// 	Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
	// 	User:     dbUser,
	// 	Password: dbPass,
	// 	Database: dbName,
	// })

	creds, err := pg.ParseURL("postgres://postgres:postgrespw@host.docker.internal:49153/postgres?sslmode=disable")
	if err != nil {
		panic("db not ready")
	}

	db := pg.Connect(creds)
	db.AddQueryHook(dbLogger{})

	satelliteRepository := pgsql.NewSatelliteRepository(db)
	topsecretService := locate.NewTopSecretService(satelliteRepository)
	srv := server.New(host, port, topsecretService)
	return srv.Run()
}
