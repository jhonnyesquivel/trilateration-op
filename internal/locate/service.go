package locate

import (
	"context"
	"math"

	quasar "github.com/jhonnyesquivel/quasar-op/internal"
)

type Point struct {
	X float64
	Y float64
}

type TopSecretService struct {
	satelliteRepository quasar.SatelliteRepository
}

func NewTopSecretService(repository quasar.SatelliteRepository) TopSecretService {
	return TopSecretService{
		satelliteRepository: repository,
	}
}

func (s TopSecretService) GetEmissorShip(ctx context.Context, satReq []quasar.Satellite) (quasar.Emissor, error) {
	var (
		satMsgs   [][]string
		distances []float64
	)

	for _, v := range satReq {
		satMsgs = append(satMsgs, v.Message().Value())
		distances = append(distances, v.Distance().Value())
	}

	message := s.getMessage(satMsgs)
	xAxis, yAxis, err := s.getLocation(ctx, distances)
	if err != nil {
		return quasar.Emissor{}, err
	}

	emissor, err := quasar.NewEmmisor(message, xAxis, yAxis)
	if err != nil {
		return quasar.Emissor{}, err
	}

	return emissor, nil
}

func (s TopSecretService) MapSatellite(name string, distance float64, msg []string) (quasar.Satellite, error) {

	satellite, err := quasar.NewSatelliteWithDistance(name, distance, msg)
	if err != nil {
		return quasar.Satellite{}, err
	}

	return satellite, nil
}

func (s TopSecretService) getMessage(messages [][]string) (msg string) {
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

// // algorithm: https://math.stackexchange.com/a/884851
func (s TopSecretService) getLocation(ctx context.Context, distances []float64) (float64, float64, error) {

	p, err := s.satelliteRepository.Fetch(ctx)
	if err != nil {
		return 0, 0, err
	}

	A := -2*p[0].Position().AxisX().Value() + 2*p[1].Position().AxisX().Value()
	B := -2*p[0].Position().AxisY().Value() + 2*p[1].Position().AxisY().Value()
	C := math.Pow(distances[0], 2) - math.Pow(distances[1], 2) - math.Pow(p[0].Position().AxisX().Value(), 2) + math.Pow(p[1].Position().AxisX().Value(), 2) - math.Pow(p[0].Position().AxisY().Value(), 2) + math.Pow(p[1].Position().AxisY().Value(), 2)

	D := -2*p[1].Position().AxisX().Value() + 2*p[2].Position().AxisX().Value()
	E := -2*p[1].Position().AxisY().Value() + 2*p[2].Position().AxisY().Value()
	F := math.Pow(distances[1], 2) - math.Pow(distances[2], 2) - math.Pow(p[1].Position().AxisX().Value(), 2) + math.Pow(p[2].Position().AxisX().Value(), 2) - math.Pow(p[1].Position().AxisY().Value(), 2) + math.Pow(p[2].Position().AxisY().Value(), 2)

	// determinants solution
	determinant := (A*E - B*D)
	X := (C*E - F*B) / determinant
	Y := (A*F - C*D) / determinant
	return X, Y, nil
}
