package locate

import (
	"context"
	"math"
	"strings"

	quasar "github.com/jhonnyesquivel/quasar-op/internal"
)

type Point struct {
	X float64
	Y float64
}

type TopSecretService struct {
	repository quasar.SatelliteRepository
}

func NewTopSecretService(repository quasar.SatelliteRepository) TopSecretService {

	return TopSecretService{
		repository: repository,
	}
}

func (s *TopSecretService) GetEmissorShip(ctx context.Context, satReq []*quasar.Satellite) (quasar.Emissor, error) {

	dbSats, err := s.repository.GetAll(ctx)
	if err != nil {
		return quasar.Emissor{}, err
	}

	satMap := make(map[string]*quasar.Satellite)
	satBase := make([]*quasar.Satellite, len(dbSats))
	for _, s := range dbSats {
		satMap[strings.ToLower(s.Name().Value())] = s
	}

	for i, s := range satReq {
		x := satMap[strings.ToLower(s.Name().Value())].X()
		y := satMap[strings.ToLower(s.Name().Value())].Y()

		sb, err := quasar.NewSatellite(s.Name().Value(), s.Distance().Value(), s.Message().Value(), x, y)
		if err != nil {
			return quasar.Emissor{}, err
		}

		satBase[i] = &sb
	}

	return s.getEmissorShip(ctx, satBase)
}

func (s *TopSecretService) GetEmissorShipFromDB(ctx context.Context) (quasar.Emissor, error) {
	satellites, err := s.repository.GetAll(ctx)
	if err != nil {
		return quasar.Emissor{}, err
	}
	return s.GetEmissorShip(ctx, satellites)
}

func (s *TopSecretService) SaveSatelliteDistance(ctx context.Context, name string, distance float64, message []string) error {

	sat, err := s.MapSatellite(name, message, distance)
	if err != nil {
		return err
	}

	err = s.repository.SaveDistance(ctx, sat)
	if err != nil {
		return err
	}

	return nil
}

func (s *TopSecretService) getEmissorShip(ctx context.Context, satReq []*quasar.Satellite) (quasar.Emissor, error) {
	var (
		satMsgs   [][]string
		distances []float64
	)

	for _, v := range satReq {
		satMsgs = append(satMsgs, v.Message().Value())
		distances = append(distances, v.Distance().Value())
	}

	message := s.getMessage(satMsgs)
	xAxis, yAxis, err := s.getLocation(satReq, distances)
	if err != nil {
		return quasar.Emissor{}, err
	}

	emissor, err := quasar.NewEmmisor(message, xAxis, yAxis)
	if err != nil {
		return quasar.Emissor{}, err
	}

	return emissor, nil
}

func (s *TopSecretService) MapSatellite(name string, msg []string, distance float64, coords ...float64) (quasar.Satellite, error) {

	satellite, err := quasar.NewSatellite(name, distance, msg, coords...)
	if err != nil {
		return quasar.Satellite{}, err
	}

	return satellite, nil
}

func (s *TopSecretService) getMessage(messages [][]string) (msg string) {
	var fullMsg []string
	var maxItms = 0
	for _, v := range messages {
		if l := len(v); l > maxItms {
			maxItms = l
		}
	}

	for j := 0; j < maxItms; j++ {
		for i := 0; i < len(messages); i++ {
			if len(messages[i]) > j {
				if word := messages[i][j]; len(word) > 0 {
					if len(fullMsg) == 0 {
						fullMsg = append(fullMsg, []string{word}...)
					} else if fullMsg[len(fullMsg)-1] != word {
						fullMsg = append(fullMsg, []string{word}...)
					}
				}
			}
		}
	}

	return strings.Join(fullMsg, " ")

}

// // algorithm: https://math.stackexchange.com/a/884851
func (s *TopSecretService) getLocation(p []*quasar.Satellite, distances []float64) (float64, float64, error) {

	A := -2*p[0].X() + 2*p[1].X()
	B := -2*p[0].Y() + 2*p[1].Y()
	C := math.Pow(distances[0], 2) - math.Pow(distances[1], 2) - math.Pow(p[0].X(), 2) + math.Pow(p[1].X(), 2) - math.Pow(p[0].Y(), 2) + math.Pow(p[1].Y(), 2)

	D := -2*p[1].X() + 2*p[2].X()
	E := -2*p[1].Y() + 2*p[2].Y()
	F := math.Pow(distances[1], 2) - math.Pow(distances[2], 2) - math.Pow(p[1].X(), 2) + math.Pow(p[2].X(), 2) - math.Pow(p[1].Y(), 2) + math.Pow(p[2].Y(), 2)

	// determinants solution
	determinant := (A*E - B*D)
	X := (C*E - F*B) / determinant
	Y := (A*F - C*D) / determinant
	return X, Y, nil
}
