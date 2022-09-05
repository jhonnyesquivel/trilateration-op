package locate

import (
	"context"
	"math"

	quasar "github.com/jhonnyesquivel/quasar-op/internal"
)

type sqlSatellite struct {
	Name     string
	Position sqlPosition
}

type sqlPosition struct {
	AxisX float64
	AxisY float64
}

type Point struct {
	X float64
	Y float64
}

var satellites = [3]sqlSatellite{
	{Name: "Kenobi", Position: sqlPosition{AxisX: -500, AxisY: -200}},
	{Name: "Skywalker", Position: sqlPosition{AxisX: 100, AxisY: -100}},
	{Name: "Sato", Position: sqlPosition{AxisX: 500, AxisY: 100}},
}

type TopSecretService struct {
}

func NewTopSecretService() TopSecretService {
	return TopSecretService{}
}

func (s TopSecretService) GetEmissorShip(ctx context.Context, satellites []quasar.Satellite) (quasar.Emissor, error) {

	var (
		satMsgs   [][]string
		distances []float64
	)

	for _, v := range satellites {
		satMsgs = append(satMsgs, v.Message().Value())
		distances = append(distances, v.Distance().Value())
	}

	message := getMessage(satMsgs)
	xAxis, yAxis := getLocation(distances)

	emissor, err := quasar.NewEmmisor(message, xAxis, yAxis)
	if err != nil {
		return quasar.Emissor{}, err
	}

	return emissor, nil
}

func (s TopSecretService) MapSatellite(name string, distance float64, msg []string) (quasar.Satellite, error) {

	satellite, err := quasar.NewSatellite(name, distance, msg)
	if err != nil {
		return quasar.Satellite{}, err
	}

	return satellite, nil
}

func getMessage(messages [][]string) (msg string) {
	it, msg := len(messages[2]), ""

	if len(messages[0]) > len(messages[1]) {
		it = len(messages[0])
	}

	if len(messages[1]) > len(messages[2]) {
		it = len(messages[1])
	}

	for i := 0; i < it; i++ {
		if len(messages[0]) > i {
			msg = msg + " " + messages[0][i]
		}
		if len(messages[1]) > i {
			msg = msg + " " + messages[1][i]
		}
		if len(messages[2]) > i {
			msg = msg + " " + messages[2][i]
		}
	}
	return msg
}

// func getLocation(distances []float64) (x, y float64) {

// 	//unit vector in a direction from point1 to point 2
// 	ksDistance :=
// 		math.Sqrt(math.Pow(float64(satellites[1].Position.AxisX)-float64(satellites[0].Position.AxisX), 2) + math.Pow(float64(satellites[1].Position.AxisY)-float64(satellites[0].Position.AxisY), 2))

// 	ex := Point{
// 		X: (float64(satellites[1].Position.AxisX) - float64(satellites[0].Position.AxisX)) / ksDistance,
// 		Y: (float64(satellites[1].Position.AxisY) - float64(satellites[0].Position.AxisY)) / ksDistance,
// 	}

// 	aux := Point{
// 		X: (float64(satellites[2].Position.AxisX) - float64(satellites[0].Position.AxisX)),
// 		Y: (float64(satellites[2].Position.AxisY) - float64(satellites[0].Position.AxisY)),
// 	}

// 	//signed magnitude of the x component
// 	i := ex.X*aux.X + ex.Y*aux.Y

// 	//the unit vector in the y direction.
// 	aux2 := Point{
// 		X: float64(satellites[2].Position.AxisX) - float64(satellites[0].Position.AxisX) - i*ex.X,
// 		Y: float64(satellites[2].Position.AxisY) - float64(satellites[0].Position.AxisY) - i*ex.X,
// 	}

// 	ey := Point{
// 		X: aux2.X / normalize(aux2),
// 		Y: aux2.Y / normalize(aux2),
// 	}

// 	//the signed magnitude of the y component
// 	j := ey.X*aux.X + ey.Y*aux.Y

// 	//coordinates
// 	x1 := (math.Pow(float64(distances[0]), 2) - math.Pow(float64(distances[1]), 2) + math.Pow(ksDistance, 2)) / (2 * ksDistance)
// 	y1 := (math.Pow(float64(distances[0]), 2)-math.Pow(float64(distances[2]), 2)+math.Pow(i, 2)+math.Pow(j, 2))/(2*j) - (i * float64(x) / j)

// 	//result coordinates
// 	x = float64(float64(satellites[0].Position.AxisX) + x1*ex.X + y1*ey.X)
// 	y = float64(float64(satellites[0].Position.AxisY) + x1*ex.Y + y1*ey.Y)

// 	return x, y
// }

func getLocation(distances []float64) (x, y float64) {
	print(distances[1])
	return trilateration(distances, satellites)

}

// algorithm: https://math.stackexchange.com/a/884851
func trilateration(r []float64, p [3]sqlSatellite) (X float64, Y float64) {

	A := -2*p[0].Position.AxisX + 2*p[1].Position.AxisX
	B := -2*p[0].Position.AxisY + 2*p[1].Position.AxisY
	C := math.Pow(r[0], 2) - math.Pow(r[1], 2) - math.Pow(p[0].Position.AxisX, 2) + math.Pow(p[1].Position.AxisX, 2) - math.Pow(p[0].Position.AxisY, 2) + math.Pow(p[1].Position.AxisY, 2)

	D := -2*p[1].Position.AxisX + 2*p[2].Position.AxisX
	E := -2*p[1].Position.AxisY + 2*p[2].Position.AxisY
	F := math.Pow(r[1], 2) - math.Pow(r[2], 2) - math.Pow(p[1].Position.AxisX, 2) + math.Pow(p[2].Position.AxisX, 2) - math.Pow(p[1].Position.AxisY, 2) + math.Pow(p[2].Position.AxisY, 2)

	// determinants solution
	determinant := (A*E - B*D)
	X = (C*E - F*B) / determinant
	Y = (A*F - C*D) / determinant
	return X, Y
}
