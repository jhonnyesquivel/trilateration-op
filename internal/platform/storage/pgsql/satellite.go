package pgsql

const (
	SqlSatelliteTable = "Satellites"
	SqlPositionTable  = "Positions"
)

type sqlSatellite struct {
	tableName  struct{} `pg:"satellites"`
	Id         string   `pg:",pk"`
	Name       string
	PositionId string
	Distance   float64
	Message    []string     `pg:",array"`
	Position   *sqlPosition `pg:",rel:has-one"`
}

type sqlPosition struct {
	tableName struct{} `pg:"positions"`
	Id        string   `pg:",pk"`
	X         float64  `pg:",notnull"`
	Y         float64  `pg:",notnull"`
}
