package topsecret

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	quasar "github.com/jhonnyesquivel/quasar-op/internal"
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

type TopSecretHandlersTestSuite struct {
	suite.Suite
	gin     *gin.Engine
	service *locate.TopSecretService
}

// this function executes before the test suite begins execution
func (suite *TopSecretHandlersTestSuite) SetupSuite() {
	sa, _ := quasar.NewSatellite("Kenobi", 670.82, []string{"este", "", "", "mensaje", ""}, -500, -200)
	sb, _ := quasar.NewSatellite("Skywalker", 200, []string{"", "es", "", "", "secreto"}, 100, -100)
	sc, _ := quasar.NewSatellite("Sato", 400, []string{"este", "", "un", "", ""}, 500, 100)
	data := []*quasar.Satellite{&sa, &sb, &sc}

	mockRepo := new(SatelliteRepositoryMock)
	mockRepo.On("GetAll", mock.Anything).Return(data, nil)

	service := locate.NewTopSecretService(mockRepo)

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.POST("/topsecret", TopSecretPOSTHandler(service))
	r.POST("/topsecret_split/:satellite", TopSecretSplitPOSTHandler(service))
	r.GET("/topsecret_split", TopSecretSplitGETHandler(service))

	suite.gin = r
	suite.service = &service
}

func (suite *TopSecretHandlersTestSuite) Test_TopSecretGETHandler_GetEmissor() {

	suite.T().Run("given an invalid request it returns 400", func(t *testing.T) {
		topsecReq := topsecretReq{
			Satelites: []topsecretSplitReq{
				{
					Name:     "kenobi",
					Distance: 0,
					Message:  []string{},
				},
			},
		}

		b, err := json.Marshal(topsecReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/topsecret", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		suite.gin.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	suite.T().Run("given a valid request it returns 200", func(t *testing.T) {
		topsecReq := topsecretReq{
			Satelites: []topsecretSplitReq{
				{
					Name:     "kenobi",
					Distance: 0,
					Message:  []string{},
				},
				{
					Name:     "Skywalker",
					Distance: 0,
					Message:  []string{},
				},
				{
					Name:     "Sato",
					Distance: 0,
					Message:  []string{},
				},
			},
		}

		b, err := json.Marshal(topsecReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/topsecret", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		suite.gin.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	suite.T().Run("given a valid request with a not findable data", func(t *testing.T) {
		topsecReq := topsecretReq{
			Satelites: []topsecretSplitReq{
				{
					Name:     "kenobi",
					Distance: 0,
					Message:  []string{},
				},
				{
					Name:     "kenobi",
					Distance: 0,
					Message:  []string{},
				},
				{
					Name:     "kenobi",
					Distance: 0,
					Message:  []string{},
				},
			},
		}

		b, err := json.Marshal(topsecReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/topsecret", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		suite.gin.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func (suite *TopSecretHandlersTestSuite) Test_TopSecretSplitPOSTHandler_SaveDistance(t *testing.T) {

	suite.T().Run("given an invalid request it returns 400", func(t *testing.T) {
		topsecReq := topsecretSplitReq{}

		b, err := json.Marshal(topsecReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/topsecret_split/:satellite", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		suite.gin.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	suite.T().Run("given a valid request it returns 200", func(t *testing.T) {
		topsecReq := topsecretSplitReq{
			Name:     "kenobi",
			Distance: 100,
			Message:  []string{},
		}

		b, err := json.Marshal(topsecReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/topsecret_split/"+topsecReq.Name, bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		suite.gin.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func (suite *TopSecretHandlersTestSuite) Test_TopSecretSplitPOSTHandler_GetEmissor(t *testing.T) {

	suite.T().Run("given a valid request it returns 200", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/topsecret_split", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		suite.gin.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

}
