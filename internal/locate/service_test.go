package locate

import (
	"context"
	"errors"
	"math"
	"testing"

	quasar "github.com/jhonnyesquivel/quasar-op/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type SatelliteRepositoryMock struct {
	mock.Mock
}

func (m SatelliteRepositoryMock) GetAll(ctx context.Context) ([]*quasar.Satellite, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*quasar.Satellite), args.Error(1)
}

func (m SatelliteRepositoryMock) SaveDistance(ctx context.Context, req quasar.Satellite) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func Test_TopSecretService_GetEmissorShipFromDB_Error(t *testing.T) {
	mockRepo := new(SatelliteRepositoryMock)
	mockRepo.On("GetAll", mock.Anything).Return([]*quasar.Satellite{}, errors.New("something unexpected happened"))

	service := NewTopSecretService(mockRepo)
	_, err := service.GetEmissorShipFromDB(context.Background())
	assert.Error(t, err)
}

func Test_TopSecretService_GetEmissorShipFromDB_Success(t *testing.T) {

	sa, err := quasar.NewSatellite("Kenobi", 670.82, []string{"este", "", "", "mensaje", ""}, -500, -200)
	require.NoError(t, err)
	sb, err := quasar.NewSatellite("Skywalker", 200, []string{"", "es", "", "", "secreto"}, 100, -100)
	require.NoError(t, err)
	sc, err := quasar.NewSatellite("Sato", 400, []string{"este", "", "un", "", ""}, 500, 100)
	require.NoError(t, err)
	data := []*quasar.Satellite{&sa, &sb, &sc}

	mockRepo := new(SatelliteRepositoryMock)
	mockRepo.On("GetAll", mock.Anything).Return(data, nil)

	service := NewTopSecretService(mockRepo)
	emissor, err := service.GetEmissorShipFromDB(context.Background())
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, "este es un mensaje secreto", emissor.Message().Value())
	assert.Equal(t, float64(100), math.RoundToEven(emissor.Position().AxisX().Value()))
	assert.Equal(t, float64(100), math.RoundToEven(emissor.Position().AxisY().Value()))
}

func Test_TopSecretService_GetEmissorShipFrom_Success(t *testing.T) {

	sa, err := quasar.NewSatellite("Kenobi", 670.82, []string{"este", "", "", "mensaje", ""}, -500, -200)
	require.NoError(t, err)
	sb, err := quasar.NewSatellite("Skywalker", 200, []string{"", "es", "", "", "secreto"}, 100, -100)
	require.NoError(t, err)
	sc, err := quasar.NewSatellite("Sato", 400, []string{"este", "", "un", "", ""}, 500, 100)
	require.NoError(t, err)
	req := []*quasar.Satellite{&sa, &sb, &sc}
	resp := []*quasar.Satellite{&sa, &sb, &sc}

	mockRepo := new(SatelliteRepositoryMock)
	mockRepo.On("GetAll", mock.Anything).Return(resp, nil)

	service := NewTopSecretService(mockRepo)

	emissor, err := service.GetEmissorShip(context.Background(), req)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetAll", 0) //check was not called the bd

	assert.Equal(t, "este es un mensaje secreto", emissor.Message().Value())
	assert.Equal(t, float64(100), math.RoundToEven(emissor.Position().AxisX().Value()))
	assert.Equal(t, float64(100), math.RoundToEven(emissor.Position().AxisY().Value()))
}

func Test_TopSecretService_SaveSatelliteDistance_Error(t *testing.T) {
	mockRepo := new(SatelliteRepositoryMock)
	mockRepo.On("SaveDistance", mock.Anything, mock.Anything).Return(errors.New("something unexpected happened"))
	service := NewTopSecretService(mockRepo)

	err := service.SaveSatelliteDistance(context.Background(), "test", 0, []string{"", ""})

	assert.Error(t, err)
}

func Test_TopSecretService_SaveSatelliteDistance_Success(t *testing.T) {

	var (
		name     = "Kenobi"
		distance = 670.82
		message  = []string{"este", "", "", "mensaje", ""}
	)

	arg, err := quasar.NewSatellite(name, distance, message)
	require.NoError(t, err)

	mockRepo := new(SatelliteRepositoryMock)
	mockRepo.On("SaveDistance", mock.Anything, arg).Return(nil)

	service := NewTopSecretService(mockRepo)

	err = service.SaveSatelliteDistance(context.Background(), name, distance, message)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "SaveDistance", 0)

	assert.NoError(t, err)
}
