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

// TopSecretGETHandler returns an HTTP handler for courses creation.
func TopSecretGETHandler(topsecretService locate.TopSecretService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			req        topsecretReq
			satellites []*quasar.Satellite
		)

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, v := range req.Satelites {
			sMap, err := topsecretService.MapSatellite(v.Name, v.Message, v.Distance)
			if err != nil {
				break
			}
			satellites = append(satellites, &sMap)
		}

		emmisor, err := topsecretService.GetEmissorShip(ctx, satellites)
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

// TopSecretHandler returns an HTTP handler for courses creation.
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

// TopSecretHandler returns an HTTP handler for courses creation.
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
