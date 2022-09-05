package topsecret

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	quasar "github.com/jhonnyesquivel/quasar-op/internal"
	"github.com/jhonnyesquivel/quasar-op/internal/locate"
)

type satelliteRequest struct {
	Satelites []struct {
		Name     string   `json:"name" binding:"required"`
		Distance float64  `json:"distance" binding:"required"`
		Message  []string `json:"message" binding:"required"`
	} `json:"satellites" binding:"required"`
}

type topsecretResponse struct {
	Position positionResponse `json:"position"`
	Message  string           `json:"message"`
}

type positionResponse struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// CreateHandler returns an HTTP handler for courses creation.
func CreateHandler(topsecretService locate.TopSecretService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			req        satelliteRequest
			satellites []quasar.Satellite
		)

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, v := range req.Satelites {
			sMap, err := topsecretService.MapSatellite(v.Name, v.Distance, v.Message)
			if err != nil {
				break
			}
			satellites = append(satellites, sMap)
		}

		emmisor, err := topsecretService.GetEmissorShip(ctx, satellites)

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
