package topsecret

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	quasar "github.com/jhonnyesquivel/quasar-op/internal"
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
)

type topsecretReq struct {
	Satelites []topsecretSplitReq `json:"satellites" binding:"required"`
}

type topsecretSplitReq struct {
	Name     string   `json:"name"`
	Distance float64  `json:"distance" binding:"required"`
	Message  []string `json:"message" binding:"required"`
}

type topsecretResponse struct {
	Position positionResponse `json:"position"`
	Message  string           `json:"message"`
}

type positionResponse struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Top Secret godoc
// @Summary Get the emissor coords
// @Description Get the emissor coords using a trilateration algorithm
// @Accept  json
// @Produce  json
// @Param satellites body topsecretReq true "Get emissor location"
// @Success 200 {object} topsecretResponse
// @Failure 400 {string} Error
// @Failure 404 {string} Error
// @Router /topsecret [post]
func TopSecretPOSTHandler(topsecretService locate.TopSecretService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			req        topsecretReq
			satellites []*quasar.Satellite
			err        error
			emmisor    quasar.Emissor
		)

		if err = ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, v := range req.Satelites {
			sMap, locErr := topsecretService.MapSatellite(v.Name, v.Message, v.Distance)
			if locErr != nil {
				err = locErr
				break
			}
			satellites = append(satellites, &sMap)
		}
		if err == nil {
			emmisor, err = topsecretService.GetEmissorShip(ctx, satellites)
		}
		if err != nil {
			switch {
			case
				errors.Is(err, quasar.ErrDistanceEmpty),
				errors.Is(err, quasar.ErrMessageEmpty),
				errors.Is(err, quasar.ErrNotEnoughSatellites),
				errors.Is(err, quasar.ErrSatelliteNotExistsOrIsMissing),
				errors.Is(err, quasar.ErrNameEmpty):
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			case
				errors.Is(err, quasar.ErrNotLocalizable):
				ctx.JSON(http.StatusNotFound, err.Error())
			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		resp := &topsecretResponse{
			Position: positionResponse{emmisor.Position().AxisX().Value(), emmisor.Position().AxisY().Value()},
			Message:  emmisor.Message().Value(),
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// Top Secret split godoc
// @Summary Save emissor distance
// @Description Save the distance from the emissor to an specific satellite
// @Accept  json
// @Param satellite body topsecretSplitReq true "Emissor distance part"
// @Param   satellite_name     path    string     true        "Satellite Name"
// @Success 200
// @Failure 400 {string} Error
// @Failure 404 {string} Error
// @Router /topsecret_split/{satellite_name} [post]
func TopSecretSplitPOSTHandler(topsecretService locate.TopSecretService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			req topsecretSplitReq
		)

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		req.Name = ctx.Param("satellite")

		err := topsecretService.SaveSatelliteDistance(ctx, req.Name, req.Distance, req.Message)

		if err != nil {
			switch {
			case errors.Is(err, quasar.ErrDistanceEmpty),
				errors.Is(err, quasar.ErrMessageEmpty), errors.Is(err, quasar.ErrNameEmpty):
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		ctx.Status(http.StatusOK)
	}
}

// Top Secret split godoc
// @Summary Save Get the emissor coords using a trilateration algorithm using stored data
// @Description Get the emissor coords using a trilateration algorithm with stored data
// @Produce  json
// @Success 200 {object} topsecretResponse
// @Failure 400 {string} Error
// @Failure 404 {string} Error
// @Router /topsecret_split [get]
func TopSecretSplitGETHandler(topsecretService locate.TopSecretService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		emmisor, err := topsecretService.GetEmissorShipFromDB(ctx)
		if err != nil {
			switch {
			case errors.Is(err, quasar.ErrDistanceEmpty),
				errors.Is(err, quasar.ErrMessageEmpty), errors.Is(err, quasar.ErrNameEmpty):
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		resp := &topsecretResponse{
			Position: positionResponse{emmisor.Position().AxisX().Value(), emmisor.Position().AxisY().Value()},
			Message:  emmisor.Message().Value(),
		}

		ctx.JSON(http.StatusOK, resp)

	}
}
