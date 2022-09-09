package pgsql

import (
	"context"
	"strings"

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

func (si *SatelliteRepository) GetAll(ctx context.Context) ([]*quasar.Satellite, error) {
	satellites := []sqlSatellite{}

	err := si.conn.Model(&satellites).Relation("Position").Limit(3).Order("order").Select()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching the satellites")
	}

	mapSat := make([]*quasar.Satellite, len(satellites))
	for i, sat := range satellites {

		tmp, err := quasar.NewSatellite(sat.Name, sat.Distance, sat.Message, sat.Position.X, sat.Position.Y)
		if err != nil {
			return nil, errors.Wrap(err, "error fetching the satellites")
		}

		mapSat[i] = &tmp
	}

	return mapSat, nil
}

func (si *SatelliteRepository) SaveDistance(ctx context.Context, req quasar.Satellite) error {
	var sqlSatellite sqlSatellite

	si.conn.Model(&sqlSatellite).Where("lower(name) = ?", strings.ToLower(req.Name().Value())).Limit(1).Select()
	if len(sqlSatellite.Id) == 0 {
		return errors.New("satellite not found")
	}

	sqlSatellite.Distance = req.Distance().Value()
	sqlSatellite.Message = req.Message().Value()

	_, err := si.conn.Model(&sqlSatellite).WherePK().Update()
	if err != nil {
		return errors.Wrap(err, "error updating satellite")
	}

	return nil
}
