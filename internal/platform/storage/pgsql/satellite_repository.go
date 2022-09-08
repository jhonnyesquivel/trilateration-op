package pgsql

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	quasar "github.com/jhonnyesquivel/quasar-op/internal"
	"github.com/jhonnyesquivel/quasar-op/internal/platform/util"
	"github.com/pkg/errors"
)

func init() {
	orm.SetTableNameInflector(func(s string) string {
		return util.ToSnakeCase(s)
	})
}

type SatelliteRepository struct {
	conn *pg.DB
}

func NewSatelliteRepository(dbConn *pg.DB) *SatelliteRepository {
	return &SatelliteRepository{
		conn: dbConn,
	}
}

func (si *SatelliteRepository) Fetch(ctx context.Context) ([]*quasar.Satellite, error) {
	satellites := []sqlSatellite{}

	err := si.conn.Model(&satellites).
		Relation("Position").
		Limit(3).
		Order("order").
		Select()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching the satellites base")
	}

	mapSat := make([]*quasar.Satellite, len(satellites))
	for i, sat := range satellites {

		tmp, err := quasar.NewSatelliteWithPosition(sat.Name, sat.Position.X, sat.Position.Y)
		if err != nil {
			return nil, errors.Wrap(err, "error fetching the satellites base")
		}

		mapSat[i] = &tmp
	}

	return mapSat, nil
}
